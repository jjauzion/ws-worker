package client

import (
	"context"
	"path"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
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

func (dh *DockerHandler) runImage(ctx context.Context, image string) error {
	dh.log.Info("running container", zap.String("image", image))
	reader, err := dh.client.ImagePull(ctx, image, types.ImagePullOptions{})
	if err != nil {
		dh.log.Error("failed to pull image", zap.Error(err))
		return err
	}
	dh.log.Info("docker says:")
	io.Copy(os.Stdout, reader)
	//buf := new(strings.Builder)
	//str := strings.Split(buf.String(), "\n")
	//for _, s := range str {
	//	dh.log.Info("docker says", zap.String("", s))
	//}

	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	volumes := []mount.Mount{
		{
			Type:   mount.TypeBind,
			Source: path.Join(dir, dh.config.WS_DOCKER_LOG_FOLDER),
			Target: "/logs",
		},
		{
			Type:   mount.TypeBind,
			Source: path.Join(dir, dh.config.WS_DOCKER_RESULT_FOLDER),
			Target: "/result",
		},
	}
	resp, err := dh.client.ContainerCreate(
		ctx,
		&container.Config{
			Image: image,
			Tty:   false,
		},
		&container.HostConfig{
			Mounts: volumes,
		},
		nil,
		nil,
		"")
	if err != nil {
		dh.log.Error("failed to create container", zap.Error(err))
		return err
	}

	err = dh.client.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{})
	if err != nil {
		dh.log.Error("failed to start container", zap.Error(err))
		return err
	}

	statusCh, errCh := dh.client.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			dh.log.Error("", zap.Error(err))
			return err
		}
	case <-statusCh:
	}

	out, err := dh.client.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		dh.log.Error("", zap.Error(err))
		return err
	}

	dh.log.Info("container says:")
	stdcopy.StdCopy(os.Stdout, os.Stderr, out)
	return nil
}
