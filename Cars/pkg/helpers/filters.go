package helpers

import (
	api_interface "cars/pkg/interface"
	"fmt"
	"strings"
	"sync"
)

type FilterResult struct {
	Passed bool
	CarID  int
}

func FilterByManufacturer(filter string, car api_interface.Models, ch chan<- FilterResult, wg *sync.WaitGroup) {
	defer wg.Done()

	maker, err := api_interface.FetchManufacturer(car.ManufacturerId)
	if err != nil {
		fmt.Println(err)
		ch <- FilterResult{Passed: false, CarID: car.Id}
		return
	}

	// if a filter value wasn't selected, or if the value matches the filter, send the result to the channel
	passed := filter == "" || maker.Name == filter
	ch <- FilterResult{Passed: passed, CarID: car.Id}
}

func FilterByCategory(filter string, car api_interface.Models, ch chan<- FilterResult, wg *sync.WaitGroup) {
	defer wg.Done()

	cat, err := api_interface.FetchCategory(car.CategoryId)
	if err != nil {
		fmt.Println(err)
		ch <- FilterResult{Passed: false, CarID: car.Id}
		return
	}

	// if a filter value wasn't selected, or if the value matches the filter, send the result to the channel
	passed := filter == "" || cat.Name == filter
	ch <- FilterResult{Passed: passed, CarID: car.Id}
}

func FilterByYear(filter int, car api_interface.Models, ch chan<- FilterResult, wg *sync.WaitGroup) {
	defer wg.Done()

	passed := filter == 0 || car.Year == filter
	ch <- FilterResult{Passed: passed, CarID: car.Id}
}

func FilterByCountry(filter string, car api_interface.Models, ch chan<- FilterResult, wg *sync.WaitGroup) {
	defer wg.Done()

	maker, err := api_interface.FetchManufacturer(car.ManufacturerId)

	if err != nil {
		fmt.Println(err)
		ch <- FilterResult{Passed: false, CarID: car.Id}
		return
	}

	passed := filter == "" || maker.Country == filter
	ch <- FilterResult{Passed: passed, CarID: car.Id}
}

func FilterByTransmission(filter string, car api_interface.Models, ch chan<- FilterResult, wg *sync.WaitGroup) {
	defer wg.Done()

	passed := filter == "" || strings.Contains(car.Specifications.Transmission, filter)
	ch <- FilterResult{Passed: passed, CarID: car.Id}
}

func FilterByDrivetrain(filter string, car api_interface.Models, ch chan<- FilterResult, wg *sync.WaitGroup) {
	defer wg.Done()

	passed := filter == "" || car.Specifications.Drivetrain == filter
	ch <- FilterResult{Passed: passed, CarID: car.Id}
}

func FilterByEngine(from, to float64, car api_interface.Models, ch chan<- FilterResult, wg *sync.WaitGroup) {
	defer wg.Done()

	engineStr := car.Specifications.Engine
	engine := engineFloat(engineStr)

	if engine == -1 {
		ch <- FilterResult{Passed: false, CarID: car.Id}
		return
	}

	passed := (from == 0 && to == 0) || (engine >= from && engine <= to)
	ch <- FilterResult{Passed: passed, CarID: car.Id}
}

func FilterByHorsepower(from, to int, car api_interface.Models, ch chan<- FilterResult, wg *sync.WaitGroup) {
	defer wg.Done()

	pwr := car.Specifications.Horsepower

	passed := (from == 0 && to == 0) || (pwr >= from && pwr <= to)
	ch <- FilterResult{Passed: passed, CarID: car.Id}
}
