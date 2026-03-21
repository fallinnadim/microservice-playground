package inbound

type LoginRequest struct {
	Email    string
	Password string
}

type LoginResponse struct {
	Token string
}

type RegisterRequest struct {
	Email    string
	Password string
}
