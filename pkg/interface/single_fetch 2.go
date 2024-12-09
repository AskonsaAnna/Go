package api_interface

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func FetchModel(id int) (Models, error) {
	data := Models{}

	url := "http://localhost:3000/api/models/" + strconv.Itoa(id)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error fetching %d: %s \n", id, err)
		return data, err
	}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		fmt.Printf("Error decoding model data for %d: %s\n", id, err)
		return data, err
	}

	return data, nil
}

func FetchManufacturer(id int) (Manufacturers, error) {
	data := Manufacturers{}

	url := "http://localhost:3000/api/manufacturers/" + strconv.Itoa(id)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error fetching %d: %s \n", id, err)
		return data, err
	}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		fmt.Printf("Error decoding manufacturer data for %d: %s\n", id, err)
		return data, err
	}

	return data, nil
}

func FetchCategory(id int) (Categories, error) {
	data := Categories{}

	url := "http://localhost:3000/api/categories/" + strconv.Itoa(id)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error fetching %d: %s \n", id, err)
		return data, err
	}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		fmt.Printf("Error decoding category data for %d: %s\n", id, err)
		return data, err
	}

	return data, nil
}
