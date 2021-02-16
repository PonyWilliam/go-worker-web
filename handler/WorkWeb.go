package handler

import (
	"context"

	log "github.com/micro/go-micro/v2/logger"

	WorkWeb "WorkWeb/proto/WorkWeb"
)

type WorkWeb struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *WorkWeb) Call(ctx context.Context, req *WorkWeb.Request, rsp *WorkWeb.Response) error {
	log.Info("Received WorkWeb.Call request")
	rsp.Msg = "Hello " + req.Name
	return nil
}

// Stream is a server side stream handler called via client.Stream or the generated client code
func (e *WorkWeb) Stream(ctx context.Context, req *WorkWeb.StreamingRequest, stream WorkWeb.WorkWeb_StreamStream) error {
	log.Infof("Received WorkWeb.Stream request with count: %d", req.Count)

	for i := 0; i < int(req.Count); i++ {
		log.Infof("Responding: %d", i)
		if err := stream.Send(&WorkWeb.StreamingResponse{
			Count: int64(i),
		}); err != nil {
			return err
		}
	}

	return nil
}

// PingPong is a bidirectional stream handler called via client.Stream or the generated client code
func (e *WorkWeb) PingPong(ctx context.Context, stream WorkWeb.WorkWeb_PingPongStream) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			return err
		}
		log.Infof("Got ping %v", req.Stroke)
		if err := stream.Send(&WorkWeb.Pong{Stroke: req.Stroke}); err != nil {
			return err
		}
	}
}
