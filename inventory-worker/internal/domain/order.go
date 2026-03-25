package domain

type OrderItem struct {
	ItemID   string `json:"itemId"`
	Name     string `json:"name"`
	Price    int    `json:"price"`
	Quantity int    `json:"quantity"`
}
