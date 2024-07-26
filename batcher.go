package batcher

import (
	"context"
	"sync"
	"time"
)

// BatchFunc is a handler type accepted by start the batcher workers
type BatchFunc func(context.Context, []any)

type payload struct {
	data any
}

// Batcher manages batches pools.
// It provides a framework to start go routine pools, push data unto a batch,
// wait for channels to drain and wait for routines to exit.
type Batcher struct {
	logger       Logger
	ch           chan payload
	closeOnce    *sync.Once
	workerMu     *sync.Mutex
	workerWG     *sync.WaitGroup
	chanWG       *sync.WaitGroup
	workers      uint
	batchSize    uint
	batchTimeout time.Duration
	closed       bool
}

// New initiates a new *Batcher.
//
// See batcher.Options for a list of options.
// Defaults: Workers = 1, BatchSize = 10,  BatchTimeout = 1s
func New(o *Options) *Batcher {
	opts := loadOptions(o)

	return &Batcher{
		logger:       opts.Logger,
		ch:           make(chan payload, opts.BatchSize),
		closeOnce:    &sync.Once{},
		workerMu:     &sync.Mutex{},
		workerWG:     &sync.WaitGroup{},
		chanWG:       &sync.WaitGroup{},
		workers:      opts.Workers,
		batchSize:    opts.BatchSize,
		batchTimeout: opts.BatchTimeout,
	}
}

// Start starts the batcher worker pools with a handler function.
// Function must follow the same signature as BatchFunc as the callback.
func (b *Batcher) Start(ctx context.Context, fn BatchFunc) {
	b.logger.Log("starting batch workers with ", b.workers, " workers")
	b.workerWG.Add(int(b.workers))
	for i := uint(0); i < b.workers; i++ {
		go worker(ctx, b, fn)
	}
}

// Add adds a data unto the batch
// Multiple concurrent calls are supported
//
// If Close() has been called, Add immediately returns an error
func (b *Batcher) Add(input any) error {
	if input == nil {
		return ErrInputIsNil
	}

	b.workerMu.Lock()
	defer b.workerMu.Unlock()
	if b.closed {
		return ErrInputChannelClosed
	}

	p := payload{
		data: input,
	}
	b.logger.Log("putting batch, payload: ", p)
	b.chanWG.Add(1)
	b.ch <- p
	return nil
}

// Close closes the input batch and
// waits for all active calls to Add to finish, then returns.
// It blocks until the batch is completely processed
func (b *Batcher) Close() {
	b.logger.Log("closing batch")
	b.closeOnce.Do(func() {
		b.workerMu.Lock()
		defer b.workerMu.Unlock()
		close(b.ch)
		b.wait()
		b.closed = true
	})
}

// wait waits for all go routines to shut down. Shutdown is triggered by calling Close
func (b *Batcher) wait() {
	b.chanWG.Wait()
	b.workerWG.Wait()
}

func worker(ctx context.Context, b *Batcher, fn BatchFunc) {
	ticker := time.NewTicker(b.batchTimeout)
	defer ticker.Stop()
	defer b.workerWG.Done()
	var batch []any

	for {
		select {
		case input, ok := <-b.ch:
			switch {
			case !ok:
				// On channel close
				b.logger.Log("the channel is closed, processing what is on the batch")
				if len(batch) > 0 {
					fn(ctx, batch)
				}
				return
			default:
				b.logger.Log("putting a new batch item")
				batch = append(batch, input.data)
				if uint(len(batch)) >= b.batchSize {
					b.logger.Log(
						"the batch is full, processing what is on the batch")
					fn(ctx, batch)
					batch = []any{}
					// reset the ticker for a start new timeout period
					ticker.Reset(b.batchTimeout)
				}
			}
			b.chanWG.Done()
		case <-ticker.C:
			b.logger.Log("the batch timeout, processing what is on the batch")
			if len(batch) > 0 {
				fn(ctx, batch)
				batch = []any{}
			}
		}
	}
}
