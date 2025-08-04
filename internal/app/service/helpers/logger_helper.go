package helpers

import (
	"erp/pkg/logger"

)

type EmitLog struct {
	Err func(err error, options ...logger.OptionLog)
}

type LoggerHelper struct {
	logger logger.Logger
}

func NewLoggerHelper(
	logger logger.Logger,
) *LoggerHelper {
	l := &LoggerHelper{
		logger: logger,
	}

	return l
}

func (h *LoggerHelper) EmitLog(operation string) EmitLog {
	return EmitLog{
		Err: func(err error, options ...logger.OptionLog) {
			h.logger.LogError(err, append(options, logger.OptionsLog.WithOperation(operation))...)
		},
	}
}

// func (h *LoggerHelper) emitErr(operation string) func(err error,options ...logger.OptionLog){
// 	return func(err error, options ...utils.OptionLog) {
// 		h.logger.LogError(err,append(options,utils.OptionsLog.WithOperation(operation))...)
// 	}
// }
