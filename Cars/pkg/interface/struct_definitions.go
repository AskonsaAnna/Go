package api_interface

type Manufacturers struct {
	Id           int    `json:"id"`
	Name         string `json:"name"`
	Country      string `json:"country"`
	FoundingYear int    `json:"foundingYear"`
} //`json:"manufacturers"`

type Categories struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
} //`json:"categories"`

type Models struct {
	Id             int    `json:"id"`
	Name           string `json:"name"`
	ManufacturerId int    `json:"manufacturerId"`
	CategoryId     int    `json:"categoryId"`
	Year           int    `json:"year"`

	Specifications struct {
		Engine       string `json:"engine"`
		Horsepower   int    `json:"horsepower"`
		Transmission string `json:"transmission"`
		Drivetrain   string `json:"drivetrain"`
	} `json:"specifications"`

	Image string `json:"image"`
} //`json:"carModels"`
