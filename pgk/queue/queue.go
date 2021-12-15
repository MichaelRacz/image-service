package queue

import (
	"context"
	"michaelracz/image-service/pgk/docker"
)

type Queue interface {
	Enqueueer
	Dequeueer
}

type Enqueueer interface {
	Enqueue(dockerFile docker.Dockerfile) bool
}

type Dequeueer interface {
	Dequeue(ctx context.Context) (docker.Dockerfile, bool)
}

// NOTE: A memory queue is used to save implemention time,
// a persistent queue is a better option.
type memoryQueue struct {
	queue chan docker.Dockerfile
}

func NewQueue(limit int) Queue {
	return memoryQueue{make(chan docker.Dockerfile, limit)}
}

func (mq memoryQueue) Enqueue(dockerfile docker.Dockerfile) bool {
	select {
	case mq.queue <- dockerfile:
		return true
	default:
		return false
	}
}

func (mq memoryQueue) Dequeue(ctx context.Context) (docker.Dockerfile, bool) {
	select {
	case df := <-mq.queue:
		return df, true
	case <-ctx.Done():
		return docker.Dockerfile(""), false
	}
}
