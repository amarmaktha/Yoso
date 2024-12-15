package models

type Product struct {
	ProductID    int     `json:"productid"`
	Name         string  `json:"name"`
	Description  string  `json:"description"`
	PriceInINR   float64 `json:"price_in_inr"`
	AvailableQty int     `json:"available_quantity"`
	IsActive     bool    `json:"is_active"`
    
}

type Look struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	PriceInINR   float64   `json:"price_in_inr"`
	AvailableQty int       `json:"available_quantity"`
	Products     []Product `json:"products"`
    Discount     float64 `json:"discount"`
}
