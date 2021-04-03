package client

import (
	"context"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"log"
	"time"

	pb "github.com/jjauzion/ws-worker/proto"
)

const (
	sleepBetweenCall = 30 * time.Second
)

func Run() {
	lg, cf, dh, err := dependencies()
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

	ctx := context.Background()
	lg.Info("pulling new task...", zap.Duration("sleep", sleepBetweenCall))
	for {
		r, err := c.StartTask(ctx, &pb.StartTaskReq{WithGPU: true})
		if getErrorCode(err) == getErrorCode(errNoTasksInQueue) {
			time.Sleep(sleepBetweenCall)
			continue
		} else if err != nil {
			lg.Error("failed to start task", zap.Error(err))
			time.Sleep(sleepBetweenCall)
			continue
		}
		lg.Info("starting task", zap.String("id", r.TaskId))
		err = dh.runImage(ctx, r.Job.DockerImage, r.Job.Env)
		if err != nil {
			lg.Error("", zap.Error(err))
		}
		_, err = c.EndTask(ctx, &pb.EndTaskReq{TaskId: r.TaskId, Error: []string{err.Error()}})
		if err != nil {
			lg.Error("failed to end task", zap.Error(err))
		}
	}
}
