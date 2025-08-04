package logger

const (
	defaultOperationName = "default"
)

type options struct {
	file       string
	method     string
	operation  string
	lineNumber int
}

type OptionLog func(o *options)

var OptionsLog options

func (_ *options) WithFilename(file string) OptionLog {
	return func(o *options) {
		o.file = file
	}
}
func (_ *options) WithMethod(file string) OptionLog {
	return func(o *options) {
		o.method = file
	}
}
func (_ *options) WithOperation(file string) OptionLog {
	return func(o *options) {
		o.operation = file
	}
}
func (_ *options) WithLineNumber(lineNumber int) OptionLog {
	return func(o *options) {
		o.lineNumber = lineNumber
	}
}

func (o *options) GetOperation() string {
	return o.operation
}

func (o *options) GetFileName() string {
	return o.file
}

func (o *options) GetMethod() string {
	return o.method
}

func (o *options) GetLineNumber() int {
	return o.lineNumber
}

func (_ *options) Apply(opts ...OptionLog) options {
	options := options{}
	for _, opt := range opts {
		opt(&options)
	}
	if options.operation == "" {
		options.operation = defaultOperationName
	}
	return options
}
