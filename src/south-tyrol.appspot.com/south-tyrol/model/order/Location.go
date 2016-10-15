package order

type Location struct {
	Address string `json:"address"`;
	Latitude float64 `json:"lat"`;
	Longitude float64 `json:"lng"`;
}
