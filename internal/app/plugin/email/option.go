package email

import "erp/pkg/logger"

type options struct {
	logger     logger.Logger
	numWorkers int
	queueSize  int
}

type Option func(c *options)

var Options options

func (options) Logger(logger logger.Logger) Option {
	return func(c *options) {
		c.logger = logger
	}
}

func (options) NumQueueSize(queueSize int) Option {
	return func(c *options) {
		c.queueSize = queueSize
	}
}

func (options) NumWorkers(numWorkers int) Option {
	return func(b *options) {
		b.numWorkers = numWorkers
	}
}

func (o options) apply(opts ...Option) options {
	ret := options{}
	for _, opt := range opts {
		opt(&ret)
	}
	if ret.numWorkers == 0 {
		ret.numWorkers = 50
	}
	if ret.queueSize == 0 {
		ret.queueSize = 1000
	}
	return ret
}
