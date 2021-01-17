// Create and maintain by Chaiyapong Lapliengtrakul (chaiyapong@3dsinteractive.com), All right reserved (2021 - Present)

package main

import (
	"sync/atomic"
	"time"

	"github.com/3dsinteractive/go-queue/queue"
)

// Batch is collection of items
type Batch struct {
	q *queue.Queue
}

// NewBatch return new buffer
func NewBatch() *Batch {
	return &Batch{
		q: queue.New(),
	}
}

// Add item into buffer
func (b *Batch) Add(item interface{}) {
	b.q.PushBack(item)
}

// Read item from buffer
func (b *Batch) Read() interface{} {
	v := b.q.PopFront()
	return v
}

// Reset clear the buffer
func (b *Batch) Reset() {
	b.q.Init()
}

// BatchFillFunc is the handler for fill the batch
type BatchFillFunc func(b *Batch, payload interface{}) error

// BatchExecFunc is the handler for execute the batch
type BatchExecFunc func(b *Batch) error

// BatchEvent is struct to manage batch life cycle event
type BatchEvent struct {
	batchSize int
	timeout   time.Duration
	fill      BatchFillFunc
	execute   BatchExecFunc
	done      BatchExecFunc
	payload   chan interface{}
	errc      chan error
}

// NewBatchEvent return new BatchEvent
func NewBatchEvent(
	batchSize int,
	timeout time.Duration,
	fill BatchFillFunc,
	execute BatchExecFunc,
	payload chan interface{},
	errc chan error,
) *BatchEvent {
	return &BatchEvent{
		batchSize: batchSize,
		timeout:   timeout,
		fill:      fill,
		execute:   execute,
		payload:   payload,
		errc:      errc,
	}
}

// Start start the batch event
// Loop will exit when payload has closed close(payload)
func (be *BatchEvent) Start() {
	fill := make(chan interface{})
	exec := make(chan bool)
	done := make(chan bool, 1)
	stop := make(chan bool, 1) // stop timer channel
	defer func() {
		close(stop)
		close(done)
		close(exec)
		close(fill)
	}()

	var n int32

	if be.timeout > 0 {
		go func(timeout time.Duration) {
			for {
				timer := time.NewTimer(timeout)
				select {
				case <-stop:
					timer.Stop()
					return
				case <-timer.C:
					i := atomic.LoadInt32(&n)
					if i > 0 {
						atomic.StoreInt32(&n, 0)
						exec <- true
					}
				}
			}
		}(be.timeout)
	}

	go func(payload chan interface{}, batchSize int, timeout time.Duration) {
		for {
			p, ok := <-payload
			if ok {
				fill <- p
				i := atomic.AddInt32(&n, 1)
				if i >= int32(batchSize) {
					atomic.StoreInt32(&n, 0)
					exec <- true
				}
			} else {
				// close everything
				if timeout > 0 {
					// exit timer
					stop <- true
				}
				// execute last batch
				exec <- true
				// exit from executor
				done <- true
				return
			}
		}
	}(be.payload, be.batchSize, be.timeout)

	batch := NewBatch()
	for {
		select {
		case payload := <-fill:
			err := be.fill(batch, payload)
			if err != nil {
				be.errc <- err
			}
		case <-exec:
			err := be.execute(batch)
			batch.Reset()
			if err != nil {
				be.errc <- err
			}
		case <-done:
			return
		}
	}
}
