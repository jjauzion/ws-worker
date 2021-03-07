package client

import (
	"context"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"log"
	"time"

	pb "github.com/jjauzion/ws-worker/proto"
)

func Run() {
	lg, cf, err := dependencies()
	if err != nil {
		log.Panic(err)
	}
	address := cf.WS_GRPC_HOST + ":" + cf.WS_GRPC_PORT
	//lg.Info("...", zap.String("address", address0))
	//address := "localhost:8090"
	lg.Info("connecting to grpc server", zap.String("address", address))
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		lg.Panic("failed to connect", zap.Error(err))
	}
	defer conn.Close()
	lg.Info("connection acquired")
	c := pb.NewApiClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.StartTask(ctx, &pb.StartTaskReq{WithGPU: false})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("image: %s\ndataset: %s\n", r.Job.DockerImage, r.Job.Dataset)
}
