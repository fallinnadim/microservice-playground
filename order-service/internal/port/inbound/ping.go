package inbound

import "context"

type PingUsecase interface {
	Ping(context.Context) (string, error)
}
