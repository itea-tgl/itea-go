package ilog

var (
	logger 	ILog
)

type ILog interface {
	Init()
	Done() bool
	Debug(v ...interface{})
	Info(v ...interface{})
	Error(v ...interface{})
	Fatal(v ...interface{})
}

type LogConstruct func() ILog

func Init(l LogConstruct, opts ...IOption) {
	if l == nil {
		l = LogConsole
	}
	logger = l()
	for _, o := range opts {
		o.do(logger)
	}
	logger.Init()
}

func Done() bool {
	return logger.Done()
}

func Debug(v ...interface{}) {
	logger.Debug(v...)
}

func Info(v ...interface{}) {
	logger.Info(v...)
}

func Error(v ...interface{}) {
	logger.Error(v...)
}

func Fatal(v ...interface{}) {
	logger.Fatal(v...)
}