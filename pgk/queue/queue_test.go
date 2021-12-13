package queue

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQueue(t *testing.T) {
	q := NewQueue(10)
	assert.True(t, q.Enqueue("str1"))
	assert.True(t, q.Enqueue("str2"))
	c := q.GetChannel()
	assert.Equal(t, "str1", <-c)
	assert.Equal(t, "str2", <-c)
}

func TestQueueLimitExceeded(t *testing.T) {
	q := NewQueue(1)
	assert.True(t, q.Enqueue("str1"))
	assert.False(t, q.Enqueue("str2"))
}

func TestQueueEnqueAfterLimitExceeded(t *testing.T) {
	q := NewQueue(1)
	assert.True(t, q.Enqueue("str1"))
	assert.False(t, q.Enqueue("str2"))
	<-q.GetChannel()
	assert.True(t, q.Enqueue("str3"))
}
