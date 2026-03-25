package request

type OrderRequest struct {
	UserId string        `json:"userId" validate:"required"`
	Items  []ItemRequest `json:"items" validate:"required,min=1"`
}

type ItemRequest struct {
	Id      string `json:"id" validate:"required"`
	Ammount int    `json:"ammount" validate:"gt=0"`
}
