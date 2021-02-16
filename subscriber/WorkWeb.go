package subscriber

import (
	"context"
	log "github.com/micro/go-micro/v2/logger"

	WorkWeb "WorkWeb/proto/WorkWeb"
)

type WorkWeb struct{}

func (e *WorkWeb) Handle(ctx context.Context, msg *WorkWeb.Message) error {
	log.Info("Handler Received message: ", msg.Say)
	return nil
}

func Handler(ctx context.Context, msg *WorkWeb.Message) error {
	log.Info("Function Received message: ", msg.Say)
	return nil
}
