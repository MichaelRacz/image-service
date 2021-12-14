package queue

import "michaelracz/image-service/pgk/docker"

type Queue interface {
	Enqueue(dockerFile docker.Dockerfile) bool
	GetChannel() <-chan docker.Dockerfile
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

//
// TODO: make symmetric, pass cancel context
//
func (mq memoryQueue) GetChannel() <-chan docker.Dockerfile {
	return mq.queue
}
