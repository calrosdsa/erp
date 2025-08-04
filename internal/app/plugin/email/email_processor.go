package email

import (
	"erp/internal/app/event-bus/event"
	"erp/internal/app/plugin/email/emailhandler"
	"erp/internal/app/plugin/email/processor"
	_logger "erp/pkg/logger"
	"erp/pkg/queue"
	"errors"
	"time"
)

type emailProcessor struct {
	queue        *queue.BoundedQueue
	numWorkers   int
	logger       _logger.Logger
	processEmail ProccessEmail
	// bytesProcessed     atomic.Uint64
	// spansProcessed     atomic.Uint64
	stopCh   chan struct{}
	emitLog  func(error, ..._logger.OptionLog)
	handlers *emailhandler.EmailHandlers
}

type queueItem struct {
	queueTime time.Time
	payload   *event.NotificationData
}

func NewEmailProccessor(handlers *emailhandler.EmailHandlers, opts ...Option) processor.EmailProcessor {
	ep := newEmailProcessor(handlers, opts...)
	ep.queue.StartConsumers(ep.numWorkers, func(item interface{}) {
		value := item.(*queueItem)
		ep.processItemFromQueue(value)
	})
	return ep
}

func newEmailProcessor(handlers *emailhandler.EmailHandlers, opts ...Option) *emailProcessor {
	options := Options.apply(opts...)
	boundedQueue := queue.NewBoundedQueue(options.queueSize, func(item interface{}) {})
	ep := emailProcessor{
		queue:      boundedQueue,
		numWorkers: options.numWorkers,
		logger:     options.logger,
		handlers:   handlers,
		emitLog: func(err error, ol ..._logger.OptionLog) {
			options.logger.LogError(err, append(ol, _logger.OptionsLog.WithOperation("email"))...)
		},
	}
	processEmailFuncs := []ProccessEmail{ep.sendEmail}
	ep.processEmail = ChainnedProccessEmail(processEmailFuncs...)
	return &ep
}

func (ep *emailProcessor) sendEmail(payload *event.NotificationData) {
	if nil == payload {
		ep.emitLog(errors.New("email payload in null"))
		return
	}
	ep.handlers.SendEmail(payload)

	//SENDING EMAIL HERE
}

func (ep *emailProcessor) ProcessEmail(payload *event.NotificationData, options processor.EmailOptions) (bool, error) {
	if ok := ep.enqueueEmail(payload); !ok {
		return false, processor.ErrBusy
	}
	return true, nil
}

func (ep *emailProcessor) enqueueEmail(payload *event.NotificationData) bool {
	item := &queueItem{
		queueTime: time.Now(),
		payload:   payload,
	}
	return ep.queue.Produce(item)
}

func (ep *emailProcessor) processItemFromQueue(item *queueItem) {
	ep.processEmail(item.payload)
}

func (ep *emailProcessor) Close() error {
	close(ep.stopCh)
	ep.queue.Stop()

	return nil
}
