package api_interface

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

func FetchModels(wg *sync.WaitGroup, ch chan<- []Models) {
	defer wg.Done()
	data := []Models{}

	url := "http://localhost:3000/api/models"

	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error fetching list of car models: %s \n", err)
		ch <- nil
		return
	}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		fmt.Printf("Error decoding data: %s\n", err)
		ch <- nil
		return
	}

	ch <- data
}

func FetchManufacturers(wg *sync.WaitGroup, ch chan<- []Manufacturers) {
	defer wg.Done()
	data := []Manufacturers{}

	url := "http://localhost:3000/api/manufacturers"

	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error fetching list of car manufacturers: %s \n", err)
		ch <- nil
		return
	}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		fmt.Printf("Error decoding data: %s\n", err)
		ch <- nil
		return
	}

	ch <- data
}

func FetchCategories(wg *sync.WaitGroup, ch chan<- []Categories) {
	defer wg.Done()
	data := []Categories{}

	url := "http://localhost:3000/api/categories"

	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error fetching list of car categories: %s \n", err)
		ch <- nil
		return
	}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		fmt.Printf("Error decoding data: %s\n", err)
		ch <- nil
		return
	}

	ch <- data
}
