package client

import (
	"google.golang.org/grpc"
	"log"
	"time"
	"context"

	pb "github.com/jjauzion/ws-worker/proto"
)

const (
	address     = "localhost:50051"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewApiClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.StartTask(ctx, &pb.StartTaskReq{WithGPU: false})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("image: %s\ndataset: %s\n", r.Job.DockerImage, r.Job.Dataset)
}