package client

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/pkg/stdcopy"
	"go.uber.org/zap"
	"io"
	"os"

	"github.com/docker/docker/client"
	"github.com/jjauzion/ws-worker/conf"
	"github.com/jjauzion/ws-worker/internal/logger"
)

type DockerHandler struct {
	client client.APIClient
	log    *logger.Logger
	config conf.Configuration
}

func (dh *DockerHandler) new(ctx context.Context, log *logger.Logger, config conf.Configuration) error {
	dh.log = log
	dh.config = config
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		dh.log.Error("", zap.Error(err))
		return err
	}
	reader, err := cli.ImagePull(ctx, "docker.io/jjauzion/ws-mock-container", types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	io.Copy(os.Stdout, reader)

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: "docker.io/jjauzion/ws-mock-container",
		//Cmd:   []string{"echo", "hello world"},
		Tty: false,
	}, nil, nil, nil, "")
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			panic(err)
		}
	case <-statusCh:
	}

	out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		panic(err)
	}

	stdcopy.StdCopy(os.Stdout, os.Stderr, out)

	dh.client = cli
	return nil
}
