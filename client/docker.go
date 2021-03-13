package client

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/jjauzion/ws-worker/conf"
	"github.com/jjauzion/ws-worker/internal/logger"
	"go.uber.org/zap"
	"io"
	"os"
)

type DockerHandler struct {
	client client.APIClient
	log    *logger.Logger
	config conf.Configuration
}

func (dh *DockerHandler) new(log *logger.Logger, config conf.Configuration) error {
	dh.log = log
	dh.config = config
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		dh.log.Error("", zap.Error(err))
		return err
	}
	dh.client = cli
	return nil
}

func (dh *DockerHandler) runImage(ctx context.Context, image string) {
	reader, err := dh.client.ImagePull(ctx, image, types.ImagePullOptions{})
	if err != nil {
		dh.log.Error("failed to pull image", zap.Error(err))
	}
	dh.log.Info("docker says:")
	io.Copy(os.Stdout, reader)
	//buf := new(strings.Builder)
	//str := strings.Split(buf.String(), "\n")
	//for _, s := range str {
	//	dh.log.Info("docker says", zap.String("", s))
	//}

	resp, err := dh.client.ContainerCreate(ctx, &container.Config{
		Image: image,
		//Cmd:   []string{"echo", "hello world"},
		Tty: false,
	}, nil, nil, nil, "")
	if err != nil {
		dh.log.Error("failed to create container", zap.Error(err))
	}

	err = dh.client.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{})
	if err != nil {
		dh.log.Error("failed to start container", zap.Error(err))
	}

	statusCh, errCh := dh.client.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			dh.log.Error("", zap.Error(err))
		}
	case <-statusCh:
	}

	out, err := dh.client.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		dh.log.Error("", zap.Error(err))
	}

	dh.log.Info("container says:")
	stdcopy.StdCopy(os.Stdout, os.Stderr, out)
}
