package request

type PaymentRequest struct {
	OrderID string
	UserID  string
	Amount  int
}
