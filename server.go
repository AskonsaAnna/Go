package main

import (
	"bytes"
	"cars/pkg/helpers"
	api_interface "cars/pkg/interface"
	"log"
	"net/http"
	"strconv"
	"sync"
	"text/template"
)

func StartServer() {
	// Serve static files from the "static" directory
	fs := http.FileServer(http.Dir("./static/"))

	// Register the handler functions
	http.HandleFunc("/", IndexHandler)
	http.HandleFunc("/fetch", FetchHandler)
	http.HandleFunc("/filter", FilterHandler)
	http.HandleFunc("/compare", compareHandler)
	http.HandleFunc("/download", downloadHandler)

	// Serve HTML templates from the "templates" directory
	http.Handle("/templates/", http.StripPrefix("/templates/", http.FileServer(http.Dir("./templates/"))))

	// Serve static files from the "static" directory
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Serve images from the "api/img" directory
	http.Handle("/api/img/", http.StripPrefix("/api/img/", http.FileServer(http.Dir("./api/img/"))))

	// Start the web server
	log.Println("Server started on http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
		return
	}
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	// Ensure the request is for the root URL
	if r.URL.Path != "/" {
		// http.NotFound(w, r)
		return
	}

	// Populate PageData with lists from the api
	once.Do(initPageData)

	updateCookies(r)

	var wg sync.WaitGroup

	wg.Add(2)

	go func() {
		defer wg.Done()
		data.UniqueYears = helpers.GetUniqueYears(data.Models)
	}()
	go func() {
		defer wg.Done()
		data.UniqueCountries = helpers.GetUniqueCountries(data.Manufacturers)
	}()

	wg.Wait()

	helpers.UpdateCarCounts(data.Result, r, w)
	data.MostFrequentModelID = findMostFrequentModelID()

	if data.Cookies == nil {
		data.Cookies = make(map[string]string)
	}

	carCountsCookie, err := r.Cookie("car_counts")
	if err != nil || carCountsCookie.Value == "" {
		data.Cookies["car_counts"] = ""
	} else {
		data.Cookies["car_counts"] = carCountsCookie.Value
	}

	// Use the global template and data
	// Execute the template with data and write to response
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// fmt.Println("Page loading time:", time.Since(*LaunchTime))
}

func updateCookies(r *http.Request) {
	cookies := make(map[string]string)
	for _, cookie := range r.Cookies() {
		cookies[cookie.Name] = cookie.Value
	}
	data.Cookies = cookies
}

func FilterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Parse the form data
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	action := r.FormValue("action")

	// Create a filter instance and populate it with form data
	filter := Filter{
		// Active:          true,
		Models:          r.FormValue("models"),
		Categories:      r.FormValue("categories"),
		Country:         r.FormValue("country"),
		Year:            parseFormInt(r, "year"),
		Transmission:    r.FormValue("transmission"),
		Drivetrain:      r.FormValue("drivetrain"),
		Engine_from:     parseFormFloat(r, "engine_from"),
		Engine_to:       parseFormFloat(r, "engine_to"),
		Horsepower_from: parseFormInt(r, "horsepower_from"),
		Horsepower_to:   parseFormInt(r, "horsepower_to"),
	}

	// Здесь вы можете реализовать логику фильтрации данных на основе параметров filter
	// Пример данных
	// Here you can implement the logic for filtering data based on the filter parameters

	if isFilterEmpty(filter) || action == "reset" {
		filter.Active = false
		data.Result = []int{}

	} else {
		filter.Active = true
		applyFilters(filter)

		helpers.UpdateCarCounts(data.Result, r, w)

	}

	data.Filters = filter

	// reset comparison
	data.Compare = false

	updateCookies(r)
	data.MostFrequentModelID = findMostFrequentModelID()

	// Пример фильтрации (замените на вашу логику)
	// Filtering example (replace with your logic)

	// Отправка результатов обратно клиенту
	// Send results back to client

	// only the models with their id in data.Result are displayed in the models-container
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// helper function for FilterHandler, parseFormInt parses an integer from form data
func parseFormInt(r *http.Request, key string) int {
	valueStr := r.FormValue(key)

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return 0 // Default to 0 if parsing fails
	}

	return value
}

func parseFormFloat(r *http.Request, key string) float64 {
	valueStr := r.FormValue(key)

	value, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		return 0.0 // Default to 0.0 if parsing fails
	}

	return value
}

func isFilterEmpty(filter Filter) bool {
	return filter.Models == "" &&
		filter.Categories == "" &&
		filter.Country == "" &&
		filter.Year == 0 &&
		filter.Transmission == "" &&
		filter.Drivetrain == "" &&
		filter.Engine_from == 0 &&
		filter.Engine_to == 0 &&
		filter.Horsepower_from == 0 &&
		filter.Horsepower_to == 0
}

func applyFilters(filter Filter) {

	result := []int{}
	ch := make(chan helpers.FilterResult, len(data.Models))
	var wg sync.WaitGroup

	for _, car := range data.Models {
		wg.Add(8)
		go helpers.FilterByManufacturer(filter.Models, car, ch, &wg)
		go helpers.FilterByCategory(filter.Categories, car, ch, &wg)
		go helpers.FilterByYear(filter.Year, car, ch, &wg)
		go helpers.FilterByCountry(filter.Country, car, ch, &wg)
		go helpers.FilterByTransmission(filter.Transmission, car, ch, &wg)
		go helpers.FilterByDrivetrain(filter.Drivetrain, car, ch, &wg)
		go helpers.FilterByEngine(filter.Engine_from, filter.Engine_to, car, ch, &wg)
		go helpers.FilterByHorsepower(filter.Horsepower_from, filter.Horsepower_to, car, ch, &wg)
	}

	// Close the channel when all goroutines have finished
	go func() {
		wg.Wait()
		close(ch)
	}()

	// Collect results from the channel

	// This part of the code is responsible for collecting and counting the results
	// from the channel for each filter applied to each car.
	filterCount := 8
	resultsMap := make(map[int]int)

	for res := range ch {
		if res.Passed {
			resultsMap[res.CarID]++
		}
	}

	// if the car passed all filters, add it's id to the result slice
	for carID, count := range resultsMap {
		if count == filterCount {
			result = append(result, carID)
		}
	}

	// data.Result is a slice of all car IDs that passed all applied filters
	data.Result = result
}

func FetchHandler(w http.ResponseWriter, r *http.Request) {
	// Parse form data from the request
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Get the model ID from the form data
	modelIDStr := r.FormValue("model_id")
	modelID, err := strconv.Atoi(modelIDStr)
	if err != nil {
		http.Error(w, "Invalid model ID", http.StatusBadRequest)
		return
	}

	// Fetch details about the car model based on the ID
	data.Car1, err = helpers.CombineData(modelID)
	if err != nil {
		http.Error(w, "Car not found", http.StatusNotFound)
		return
	}

	// Parse the template file
	tmplt, err := template.ParseFiles("./templates/details.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Execute the template with the data and write the result to the response writer
	err = tmplt.Execute(w, data)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func compareHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	car1ID, err1 := strconv.Atoi(r.FormValue("car1"))
	car2ID, err2 := strconv.Atoi(r.FormValue("car2"))

	if err1 != nil || err2 != nil {
		http.Error(w, "Invalid car IDs", http.StatusBadRequest)
		return
	}

	var car1, car2 *api_interface.Models
	for _, car := range data.Models {
		if car.Id == car1ID {
			car1 = &car
		}
		if car.Id == car2ID {
			car2 = &car
		}
	}

	if car1 == nil || car2 == nil {
		http.Error(w, "Car not found", http.StatusNotFound)
		return
	}

	// reset filter
	data.Result = []int{}
	data.Filters = Filter{}

	updateCookies(r)

	// example usage of CombineData
	var wg sync.WaitGroup

	wg.Add(2)
	go func() {
		defer wg.Done()

		data.Car1, err = helpers.CombineData(car1ID)
		if err != nil {
			http.Error(w, "Car not found", http.StatusNotFound)
			return
		}

		// fmt.Println(data.Car1)
	}()

	go func() {
		defer wg.Done()

		data.Car2, err = helpers.CombineData(car2ID)
		if err != nil {
			http.Error(w, "Car not found", http.StatusNotFound)
			return
		}

		// fmt.Println(data.Car2)
	}()

	wg.Wait()

	data.Diff.Hp = helpers.PowerDelta(data.Car1.Specifications.Horsepower, data.Car2.Specifications.Horsepower)
	data.Diff.Engine = helpers.EngineDelta(data.Car1.Specifications.Engine, data.Car2.Specifications.Engine)

	data.Compare = true

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	// Generate the table text
	text := generateTableText()

	// Set the appropriate headers
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Disposition", "attachment; filename=table.txt")

	// Write the text to the response
	w.Write([]byte(text))
}

func generateTableText() string {
	var buf bytes.Buffer

	// Create a helper function for padding
	pad := func(input string, length int) string {
		for len(input) < length {
			input += " "
		}
		return input
	}

	buf.WriteString("---------------------------------------------------\n")
	buf.WriteString(pad("Characteristic", 15) + "| " + pad(data.Car1.Name, 17) + "| " + data.Car2.Name + "\n")
	buf.WriteString("---------------------------------------------------\n")
	buf.WriteString(pad("Manufacturer", 15) + "| " + pad(data.Car1.Manufacturer.Name, 17) + "| " + data.Car2.Manufacturer.Name + "\n")
	buf.WriteString(pad("Country", 15) + "| " + pad(data.Car1.Manufacturer.Country, 17) + "| " + data.Car2.Manufacturer.Country + "\n")
	buf.WriteString(pad("Year", 15) + "| " + pad(strconv.Itoa(data.Car1.Year), 17) + "| " + strconv.Itoa(data.Car2.Year) + "\n")
	buf.WriteString(pad("Engine", 15) + "| " + pad(data.Car1.Specifications.Engine, 17) + "| " + data.Car2.Specifications.Engine + "\n")
	buf.WriteString(pad("Horsepower", 15) + "| " + pad(strconv.Itoa(data.Car1.Specifications.Horsepower), 17) + "| " + strconv.Itoa(data.Car2.Specifications.Horsepower) + "\n")
	buf.WriteString(pad("Transmission", 15) + "| " + pad(data.Car1.Specifications.Transmission, 17) + "| " + data.Car2.Specifications.Transmission + "\n")
	buf.WriteString(pad("Drivetrain", 15) + "| " + pad(data.Car1.Specifications.Drivetrain, 17) + "| " + data.Car2.Specifications.Drivetrain + "\n")

	return buf.String()
}
