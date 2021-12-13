package queue

type Queue interface {
	Enqueue(string) bool
	GetChannel() <-chan string
}

type memoryQueue struct {
	queue chan string
}

func NewQueue(limit int) Queue {
	return memoryQueue{make(chan string, limit)}
}

func (mq memoryQueue) Enqueue(str string) bool {
	select {
	case mq.queue <- str:
		return true
	default:
		return false
	}
}

func (mq memoryQueue) GetChannel() <-chan string {
	return mq.queue
}
