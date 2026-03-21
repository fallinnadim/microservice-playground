package response

import "time"

type SuccessResponse[T any] struct {
	Timestamp time.Time `json:"timestamp"`
	Data      T         `json:"data,omitempty"`
}
