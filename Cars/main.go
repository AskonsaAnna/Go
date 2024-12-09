package main

import (
	"cars/pkg/helpers"
	api_interface "cars/pkg/interface"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"text/template"
)

// Global template and data to avoid parsing multiple times
var (
	tmpl *template.Template
	data *PageData
	once sync.Once
	// LaunchTime *time.Time
)

type PageData struct {
	Result  []int
	Filters Filter
	Compare bool

	Car1 helpers.Combined
	Car2 helpers.Combined
	Diff struct {
		Engine string
		Hp     string
	}

	Models        []api_interface.Models
	Manufacturers []api_interface.Manufacturers
	Categories    []api_interface.Categories

	UniqueYears     []int
	UniqueCountries []string

	Cookies             map[string]string
	MostFrequentModelID int
}

type Filter struct {
	Active          bool
	Models          string
	Categories      string
	Year            int
	Country         string
	Transmission    string
	Drivetrain      string
	Engine_from     float64
	Engine_to       float64
	Horsepower_from int
	Horsepower_to   int
}

// Define a custom template function named "in"
func in(value int, slice []int) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

func matches() int {
	return len(data.Result)
}

func previouseng(value float64) string {
	if value != 0 {
		return fmt.Sprintf("(%.1fL)", value)
	} else {
		return ""
	}
}

func previouspwr(value int) string {
	if value != 0 {
		return fmt.Sprintf("(%dHp)", value)
	} else {
		return ""
	}
}

func init() {
	// now := time.Now()
	// LaunchTime = &now

	var err error
	// tmpl, err = template.ParseFiles("./templates/index.html")

	tmpl, err = template.New("index.html").Funcs(template.FuncMap{
		"in":                in,
		"matches":           matches,
		"previouse":         previouseng,
		"previousp":         previouspwr,
		"parseCarCounts":    helpers.ParseCarCounts,
		"generateTableText": generateTableText,
	}).ParseFiles("./templates/index.html")

	if err != nil {
		log.Fatal(err)
	}
}

func initPageData() {
	// start := time.Now() // Record the start time

	var wg sync.WaitGroup
	wg.Add(3)

	// Initialize PageData concurrently
	data = &PageData{}

	// Create channels
	modelsCh := make(chan []api_interface.Models)
	manufacturersCh := make(chan []api_interface.Manufacturers)
	categoriesCh := make(chan []api_interface.Categories)

	// Fetch Models
	go api_interface.FetchModels(&wg, modelsCh)

	// Fetch Manufacturers
	go api_interface.FetchManufacturers(&wg, manufacturersCh)

	// Fetch Categories
	go api_interface.FetchCategories(&wg, categoriesCh)

	// Wait for all fetch operations to complete
	go func() {
		wg.Wait()
		close(modelsCh)
		close(manufacturersCh)
		close(categoriesCh)
	}()

	// Collect the data from channels
	data.Models = <-modelsCh
	data.Manufacturers = <-manufacturersCh
	data.Categories = <-categoriesCh

	// elapsed := time.Since(start)
	// fmt.Println("initPageData execution time:", elapsed)
}

func findMostFrequentModelID() int {
	modelCounts := make(map[int]int)

	// Extract and parse the car_counts string from the Cookies map
	carCounts, ok := data.Cookies["car_counts"]
	if !ok {
		return 1 // Return 0 if the car_counts key doesn't exist
	}

	// Split the carCounts string by '|'
	pairs := strings.Split(carCounts, "|")
	for _, pair := range pairs {
		parts := strings.Split(pair, ":")
		if len(parts) != 2 {
			continue
		}

		modelID, err := strconv.Atoi(parts[0])
		if err != nil {
			continue
		}
		count, err := strconv.Atoi(parts[1])
		if err != nil {
			continue
		}

		modelCounts[modelID] += count
	}

	// Determine the model ID with the highest count
	var mostFrequentID int
	var highestCount int
	first := true

	for id, count := range modelCounts {
		if first || count > highestCount || (count == highestCount && id < mostFrequentID) {
			mostFrequentID = id
			highestCount = count
			first = false
		}
	}

	return mostFrequentID
}

func main() {
	StartServer()

}
