package batcher

import "time"

// Options are options to pass to New()
type Options struct {
	Logger       Logger
	Workers      uint
	BatchSize    uint
	BatchTimeout time.Duration
}

func loadOptions(o *Options) *Options {
	opts := &Options{
		Workers:      1,
		BatchTimeout: time.Second,
		BatchSize:    10,
		Logger:       newDefaultLogger(),
	}
	if o == nil {
		return opts
	}

	if o.Logger != nil {
		opts.Logger = o.Logger
	}

	if o.BatchTimeout > 0 {
		opts.BatchTimeout = o.BatchTimeout
	}

	if o.BatchSize > 0 {
		opts.BatchSize = o.BatchSize
	}

	if o.Workers > 0 {
		opts.Workers = o.Workers
	}

	return opts
}
