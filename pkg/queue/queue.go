package queue

import (
	"context"
	"michaelracz/image-service/pkg/docker"
)

// Queue provides a queue of Dockerfiles
type Queue interface {
	Enqueueer
	Dequeueer
}

// Enqueueer adds Dockerfiles to a queue
type Enqueueer interface {
	Enqueue(dockerFile docker.Dockerfile) bool
}

// Dequeueer fetches Dockerfiles from a queue
type Dequeueer interface {
	Dequeue(ctx context.Context) (docker.Dockerfile, bool)
}

// NOTE: A memory queue is used to save implemention time,
// a persistent queue is a better option.
type memoryQueue struct {
	queue chan docker.Dockerfile
}

// NewQueue initializes an in memory queue
func NewQueue(limit int) Queue {
	return memoryQueue{make(chan docker.Dockerfile, limit)}
}

// Enqueue adds Dockerfiles to a queue
func (mq memoryQueue) Enqueue(dockerfile docker.Dockerfile) bool {
	select {
	case mq.queue <- dockerfile:
		return true
	default:
		return false
	}
}

// Dequeue fetches Dockerfiles from a queue
func (mq memoryQueue) Dequeue(ctx context.Context) (docker.Dockerfile, bool) {
	select {
	case df := <-mq.queue:
		return df, true
	case <-ctx.Done():
		return docker.Dockerfile(""), false
	}
}
