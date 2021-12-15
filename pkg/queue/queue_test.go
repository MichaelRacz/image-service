package queue

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQueue(t *testing.T) {
	ctx := context.Background()
	q := NewQueue(10)
	assert.True(t, q.Enqueue("dockerfile1"))
	assert.True(t, q.Enqueue("dockerfile2"))
	dequeued, _ := q.Dequeue(ctx)
	assert.Equal(t, "dockerfile1", dequeued.String())
	dequeued, _ = q.Dequeue(ctx)
	assert.Equal(t, "dockerfile2", dequeued.String())
}

func TestQueueLimitExceeded(t *testing.T) {
	q := NewQueue(1)
	assert.True(t, q.Enqueue("dockerfile1"))
	assert.False(t, q.Enqueue("dockerfile2"))
}

func TestQueueEnqueAfterLimitExceeded(t *testing.T) {
	ctx := context.Background()
	q := NewQueue(1)
	assert.True(t, q.Enqueue("dockerfile1"))
	assert.False(t, q.Enqueue("dockerfile2"))
	_, _ = q.Dequeue(ctx)
	assert.True(t, q.Enqueue("dockerfile3"))
}

func TestDequeueStopsWhenContextCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	q := NewQueue(1)
	go cancel()
	_, ok := q.Dequeue(ctx)
	assert.False(t, ok)
}
