package dispatch

import (
	"context"
	"errors"
	"michaelracz/image-service/pkg/docker"
	"michaelracz/image-service/pkg/queue"
	"testing"

	"github.com/stretchr/testify/assert"
)

type dockerClientMock struct {
	t                  *testing.T
	expectedCtx        context.Context
	expectedDockerfile docker.Dockerfile
	buildCancelFunc    context.CancelFunc
	pushCancelFunc     context.CancelFunc
	buildError         error
	pushError          error
	tagUsedInBuild     **string
}

func (dc dockerClientMock) BuildImage(ctx context.Context,
	dockerFile docker.Dockerfile, tag string) error {

	if dc.buildCancelFunc != nil {
		defer dc.buildCancelFunc()
	}

	if dc.buildError != nil {
		return dc.buildError
	}

	assert.Same(dc.t, dc.expectedCtx, ctx)
	assert.Equal(dc.t, dc.expectedDockerfile, dockerFile)
	assert.Equal(dc.t, "tag", tag)
	*dc.tagUsedInBuild = &tag
	return nil
}

func (dc dockerClientMock) PushImage(ctx context.Context, tag string) error {
	if dc.pushCancelFunc != nil {
		defer dc.pushCancelFunc()
	}

	if dc.pushError != nil {
		return dc.pushError
	}

	assert.Same(dc.t, dc.expectedCtx, ctx)
	assert.Equal(dc.t, **dc.tagUsedInBuild, tag)
	return nil
}

func (dc dockerClientMock) CreateTag() string {
	return "tag"
}

func TestDispatch(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	q := queue.NewQueue(1)
	df := docker.Dockerfile("dummy")
	dcm := dockerClientMock{t, ctx, df, nil, cancel, nil, nil, new(*string)}

	q.Enqueue(df)
	Dispatch(ctx, q, dcm)
}

func TestDispatchCompensatesBuildError(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	q := queue.NewQueue(1)
	df := docker.Dockerfile("dummy")
	dcm := dockerClientMock{t, ctx, df, cancel, nil, errors.New("build error"), nil, new(*string)}

	q.Enqueue(df)
	Dispatch(ctx, q, dcm)
}

func TestDispatchCompensatesPushError(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	q := queue.NewQueue(1)
	df := docker.Dockerfile("dummy")
	dcm := dockerClientMock{t, ctx, df, nil, cancel, nil, errors.New("push error"), new(*string)}

	q.Enqueue(df)
	Dispatch(ctx, q, dcm)
}
