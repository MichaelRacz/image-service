package queue

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQueue(t *testing.T) {
	q := NewQueue(10)
	assert.True(t, q.Enqueue("dockerfile1"))
	assert.True(t, q.Enqueue("dockerfile2"))
	c := q.GetChannel()
	assert.Equal(t, "dockerfile1", (<-c).String())
	assert.Equal(t, "dockerfile2", (<-c).String())
}

func TestQueueLimitExceeded(t *testing.T) {
	q := NewQueue(1)
	assert.True(t, q.Enqueue("dockerfile1"))
	assert.False(t, q.Enqueue("dockerfile2"))
}

func TestQueueEnqueAfterLimitExceeded(t *testing.T) {
	q := NewQueue(1)
	assert.True(t, q.Enqueue("dockerfile1"))
	assert.False(t, q.Enqueue("dockerfile2"))
	<-q.GetChannel()
	assert.True(t, q.Enqueue("dockerfile3"))
}
