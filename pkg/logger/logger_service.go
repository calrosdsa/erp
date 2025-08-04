package logger

import (
	"os"

	"github.com/spf13/viper"

	"context"

	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/sdk/log"

	l "go.opentelemetry.io/otel/log"

	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"

	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/bridges/otellogrus"
)

type loggerService struct {
	logger         *logrus.Logger
}

func New(args ...interface{}) Logger {
	var (
		serviceName string
		serviceVercion string
	)
	if len(args) == 1 {
		serviceName = args[0].(string)
	}
	if len(args) == 2 {
		serviceVercion = args[1].(string)
	}
	ctx := context.Background()
	// Create resource.
	res, err := newResource(serviceName, serviceVercion)
	if err != nil {
		panic(err)
	}

	logger := &logrus.Logger{
		Out:       os.Stderr,
		Formatter: new(logrus.TextFormatter),
		Hooks:     make(logrus.LevelHooks),
		Level:     logrus.Level(l.SeverityFatal4),
	}

	// Create a logger provider.
	// You can pass this instance directly when creating bridges.
	loggerProvider, err := newLoggerProvider(ctx, res)
	if err != nil {
		panic(err)
	}
	// Handle shutdown properly so nothing leaks.
	// defer func() {
	// 	if err := loggerProvider.Shutdown(ctx); err != nil {
	// 		fmt.Println(err)
	// 	}
	// }()
	// Register as global logger provider so that it can be accessed global.LoggerProvider.
	// Most log bridges use the global logger provider as default.
	// If the global logger provider is not set then a no-op implementation
	// is used, which fails to generate data.
	global.SetLoggerProvider(loggerProvider)

	// Create an *otellogrus.Hook and use it in your application.
	hook := otellogrus.NewHook(
		"logger",
		otellogrus.WithLoggerProvider(loggerProvider),
		otellogrus.WithLevels(AllawordLevels),
	)
	// Set the newly created hook as a global logrus hook
	logger.AddHook(hook)

	return &loggerService{
		logger:      logger,
		// serviceName: serviceName,
	}
}


func (s *loggerService) EmitLog(operation string)EmitLog  {
	return EmitLog{
		Err: func(err error, options ...OptionLog) {
			// s.logger.LogError(err, append(options, logger.OptionsLog.WithOperation(operation))...)
			s.LogError(err,append(options,OptionsLog.WithOperation(operation))...)
		},
	}
}

func (s *loggerService) LogError(err error, opts ...OptionLog) {
	s.Log(err.Error(), int(l.SeverityError), opts...)
}

func (s *loggerService) LogInfo(m string, opts ...OptionLog) {
	s.Log(m, int(l.SeverityInfo), opts...)
}

func (s *loggerService) LogWarn(m string, opts ...OptionLog) {
	s.Log(m, int(l.SeverityWarn), opts...)
}

func (s *loggerService) Log(m string, severity int, opts ...OptionLog) {
	options := OptionsLog.Apply(opts...)
	fields := logrus.Fields{}
	if options.GetFileName() != "" {
		fields["file"] = options.GetFileName()
	}
	if options.GetMethod() != "" {
		fields["method"] = options.GetMethod()
	}
	if options.GetLineNumber() == 0 {
		fields["lineNumber"]  = options.GetLineNumber()
	}
	fields["operation"] = options.GetOperation()
	entry := s.logger.WithFields(fields)
	entry.Log(logrus.Level(severity), m)
}

func newResource(serviceName, serviceVersion string) (*resource.Resource, error) {
	return resource.Merge(resource.Default(),
		resource.NewWithAttributes(semconv.SchemaURL,
			semconv.ServiceName(serviceName),
			semconv.ServiceVersion(serviceVersion),
		))
}

func newLoggerProvider(ctx context.Context, res *resource.Resource) (*log.LoggerProvider, error) {
	loggerEndpoint := viper.GetString("logger.endpoint")
	if loggerEndpoint == "" {
		loggerEndpoint = "localhost:8000"
	}
	exporter, err := otlploghttp.New(ctx,
		otlploghttp.WithEndpoint("localhost:8000"),
		otlploghttp.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}
	processor := log.NewBatchProcessor(exporter)
	provider := log.NewLoggerProvider(
		log.WithResource(res),
		log.WithProcessor(processor),
	)
	return provider, nil
}

