package domain

type Order struct {
	ID     string
	UserId string
	Status string
	Items  []OrderItem
}

type OrderItem struct {
	ItemID   string `json:"itemId"`
	Name     string `json:"name"`
	Price    int    `json:"price"`
	Quantity int    `json:"quantity"`
}
