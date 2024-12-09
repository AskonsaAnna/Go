package helpers

import (
	api_interface "cars/pkg/interface"
	"sort"
)

func GetUniqueYears(models []api_interface.Models) []int {
	seenYears := make(map[int]bool)
	var uniqueYears []int

	for _, model := range models {
		if !seenYears[model.Year] {
			seenYears[model.Year] = true
			uniqueYears = append(uniqueYears, model.Year)
		}
	}

	sort.Ints(uniqueYears)

	return uniqueYears
}

func GetUniqueCountries(manufacturers []api_interface.Manufacturers) []string {
	seenCountries := make(map[string]bool)
	var uniqueCountries []string

	for _, manufacturer := range manufacturers {
		if !seenCountries[manufacturer.Country] {
			seenCountries[manufacturer.Country] = true
			uniqueCountries = append(uniqueCountries, manufacturer.Country)
		}
	}

	return uniqueCountries
}
