package helpers

import (
	api_interface "cars/pkg/interface"
)

type Combined struct {
	Id   int
	Name string

	Manufacturer struct {
		// Id           int
		Name         string
		Country      string
		FoundingYear int
	}

	Category string
	Year     int

	Specifications struct {
		Engine       string
		Horsepower   int
		Transmission string
		Drivetrain   string
	}

	Image string
}

func CombineData(id int) (Combined, error) {

	model, err := api_interface.FetchModel(id)
	if err != nil {
		return Combined{}, err
	}

	manufacturer, err := api_interface.FetchManufacturer(model.ManufacturerId)
	if err != nil {
		return Combined{}, err
	}

	category, err := api_interface.FetchCategory(model.CategoryId)
	if err != nil {
		return Combined{}, err
	}

	combined := Combined{
		Id:   model.Id,
		Name: model.Name,

		Manufacturer: struct {
			// Id           int
			Name         string
			Country      string
			FoundingYear int
		}{
			// Id:           manufacturer.Id,
			Name:         manufacturer.Name,
			Country:      manufacturer.Country,
			FoundingYear: manufacturer.FoundingYear,
		},

		Category: category.Name,
		Year:     model.Year,

		Specifications: struct {
			Engine       string
			Horsepower   int
			Transmission string
			Drivetrain   string
		}{
			Engine:       model.Specifications.Engine,
			Horsepower:   model.Specifications.Horsepower,
			Transmission: model.Specifications.Transmission,
			Drivetrain:   model.Specifications.Drivetrain,
		},

		Image: model.Image,
	}

	return combined, nil
}
