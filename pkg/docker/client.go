package docker

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
)

const IMAGE_PUSH_TIMEOUT = time.Second * 120

type Client interface {
	BuildImage(ctx context.Context, dockerFile Dockerfile, tag string) error
	PushImage(ctx context.Context, tag string) error
	CreateTag() string
}

type dockerClient struct {
	cli               *client.Client
	authConfigEncoded *string
	registryUserID    *string
}

func NewClient(registryUserID string, password string, registryServerAddress string) (Client, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}

	authConfigEncoded, err := encodeAuthConfig(registryUserID, password, registryServerAddress)
	if err != nil {
		return nil, err
	}

	return dockerClient{cli, authConfigEncoded, &registryUserID}, nil
}

func encodeAuthConfig(registryUserID string, password string, registryServerAddress string) (*string, error) {
	var authConfig = types.AuthConfig{
		Username:      registryUserID,
		Password:      password,
		ServerAddress: registryServerAddress,
	}

	authConfigBytes, err := json.Marshal(authConfig)
	if err != nil {
		return nil, err
	}

	authConfigEncoded := base64.URLEncoding.EncodeToString(authConfigBytes)
	return &authConfigEncoded, nil
}

func (dc dockerClient) BuildImage(ctx context.Context, dockerFile Dockerfile, tag string) error {
	archive, err := archive.Generate("Dockerfile", dockerFile.String())
	if err != nil {
		return err
	}

	options := types.ImageBuildOptions{
		Tags: []string{tag},
	}

	res, err := dc.cli.ImageBuild(ctx, archive, options)
	if err != nil {
		return err
	}

	// NOTE: Printing to stdout for convenience, would be too verbose in prod
	io.Copy(os.Stdout, res.Body)
	defer res.Body.Close()

	return nil
}

func (dc dockerClient) PushImage(ctx context.Context, tag string) error {
	ctx, cancel := context.WithTimeout(ctx, IMAGE_PUSH_TIMEOUT)
	defer cancel()

	opts := types.ImagePushOptions{
		RegistryAuth: *dc.authConfigEncoded,
	}

	res, err := dc.cli.ImagePush(ctx, tag, opts)
	if err != nil {
		return err
	}

	// NOTE: Printing to stdout for convenience, would be too verbose in prod
	io.Copy(os.Stdout, res)
	defer res.Close()

	return nil
}

func (dc dockerClient) CreateTag() string {
	// NOTE: Meaningful tag naming required for prod use
	tag := fmt.Sprintf("%s/image-service-output-%d", *dc.registryUserID, time.Now().UnixMilli())
	return tag
}
