package order

type Product struct {
	MaxPrice   float64 `json:"max_price"`
	Name       string `json:"name"`
	Notes      string `json:"notes"`
	Quantity   int `json:"quantity"`
	Weight     float64 `json:"weight"`
	WeightUnit string `json:"weight_unit"`
}
