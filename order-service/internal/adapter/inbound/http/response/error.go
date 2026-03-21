package response

import "time"

type ErrorResponse struct {
	Timestamp time.Time `json:"timestamp"`
	Error     string    `json:"error"`
}
