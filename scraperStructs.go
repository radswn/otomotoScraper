package main

type InputUrl struct {
	Url  string `json:"url"`
	Page int    `json:"page"`
}

type Offer struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Url         string `json:"url"`
	JpgUrl      string `json:"jpg_url"`
	Price       int    `json:"price"`
	Year        int    `json:"year"`
	Mileage     int    `json:"mileage"`
	Engine      int    `json:"engine"`
	EngineType  string `json:"engine_type"`
	Location    string `json:"location"`
}
