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
	lg.Info("connecting to grpc server", zap.String("address", address))
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		lg.Panic("failed to connect", zap.Error(err))
	}
	defer conn.Close()
	lg.Info("connection acquired")
	c := pb.NewApiClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 45 * time.Second)
	defer cancel()
	for {
		time.Sleep(5 * time.Second)
		r, err := c.StartTask(ctx, &pb.StartTaskReq{WithGPU: false})
		if getErrorCode(err) == getErrorCode(errNoTasksInQueue) {
			lg.Info("no task in queue")
			continue
		} else if err != nil {
			lg.Error("failed to start task", zap.Error(err))
			return
		}
		lg.Info("start task image", zap.String("image", r.Job.DockerImage), zap.String("dataset", r.Job.Dataset))
	}
}
