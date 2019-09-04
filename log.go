package glog

import (
	"io/ioutil"
	"reflect"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	// LevelPanic panice log level
	LevelPanic = "panic"
	// LevelFatal fatal log level
	LevelFatal = "fatal"
	// LevelError error log level
	LevelError = "error"
	// LevelWarn warn log level
	LevelWarn = "warn"
	// LevelInfo info log level
	LevelInfo = "info"
	// LevelDebug debug log level
	LevelDebug = "debug"
	// LevelTrace trace log level
	LevelTrace = "trace"

	// DefaultLogLevel default application log level
	DefaultLogLevel = LevelTrace
)

type (
	Formatter func(input interface{}, format string) string

	// Config interface
	Config interface {
		Level() string
		Formatters() map[reflect.Type]Formatter
	}

	// DefaultConfig default config
	DefaultConfig struct {
		Config
		level        string
		formatterMap map[reflect.Type]Formatter
	}

	// Logger logger
	Logger struct {
		logrus.Entry
	}
)

var (
	defaultFormatterMap = map[reflect.Type]Formatter{
		reflect.TypeOf(time.Time{}):    ToISOTime,
		reflect.TypeOf(reflect.Struct): ToISOTime,
	}

	config Config = &DefaultConfig{
		level:        DefaultLogLevel,
		formatterMap: defaultFormatterMap,
	}
)

// Initialize logger
func Initialize(cfg Config) {
	config = cfg
	level := getLogLevel(config.Level())
	logrus.SetLevel(level)
}

// Level return log level
func (cfg *DefaultConfig) Level() string {
	return cfg.level
}

// SetLevel set log level
func (cfg *DefaultConfig) SetLevel(level string) {
	cfg.level = level
}

// Formatters return formaters
func (cfg *DefaultConfig) Formatters() map[reflect.Type]Formatter {
	return cfg.formatterMap
}

// SetFormatter set formatter
func (cfg *DefaultConfig) SetFormatter(t reflect.Type, formatter Formatter) {
	if cfg.formatterMap == nil {
		panic("Formatters was not initialized")
	}
	cfg.formatterMap[t] = formatter
}

// GetRoot get logger
func GetRoot() *Logger {
	return getLogger("", DefaultLogLevel)
}

// GetLogger get logger
func GetLogger(module string) *Logger {
	cfg := getConfig(module)
	level := DefaultLogLevel
	if cfg != nil {
		level = cfg.Level()
	}
	logger := getLogger(module, level)
	if cfg == nil { // If module is not configured, disable logger
		logger.Logger.Out = ioutil.Discard
	}
	return logger
}

// GetLoggerByPackage get logger by package name
func GetLoggerByPackage(pkg interface{}) *Logger {
	module := getSimplePackageName(pkg)
	cfg := getConfig(module)
	return getLogger(module, cfg.Level())
}

// IsLevel log
func (logger *Logger) IsLevel(level string) bool {
	return logger.Level == getLogLevel(level)
}

// Tracef log
func (logger *Logger) Tracef(format string, args ...interface{}) {
	logger.Logf(logrus.TraceLevel, format, args...)
}

// Debugf log
func (logger *Logger) Debugf(format string, args ...interface{}) {
	logger.Logf(logrus.DebugLevel, format, args...)
}

// DebugWithErrorf log
func (logger *Logger) DebugWithErrorf(format string, err error, args ...interface{}) {
	logger.LogWithErrorf(logrus.DebugLevel, format, err, args...)
}

// Infof log
func (logger *Logger) Infof(format string, args ...interface{}) {
	logger.Logf(logrus.InfoLevel, format, args...)
}

// Warnf log
func (logger *Logger) Warnf(format string, err error, args ...interface{}) {
	logger.LogWithErrorf(logrus.WarnLevel, format, err, args...)
}

// Errorf log
func (logger *Logger) Errorf(format string, err error, args ...interface{}) {
	logger.LogWithErrorf(logrus.ErrorLevel, format, err, args...)
}

// Fatalf log
func (logger *Logger) Fatalf(format string, err error, args ...interface{}) {
	logger.LogWithErrorf(logrus.FatalLevel, format, err, args...)
	logger.Logger.Exit(1)
}

// LogWithErrorf log with error
func (logger *Logger) LogWithErrorf(level logrus.Level, format string, err error, args ...interface{}) {
	if args == nil {
		args = make([]interface{}, 0)
	}
	if err != nil {
		format += "\nError: %s"
		args = append(args, AsErrStrackTrace(err))
	}
	logger.Logf(level, format, args...)
}

// Logf log
func (logger *Logger) Logf(level logrus.Level, format string, args ...interface{}) {
	if len(args) > 0 {
		args = logger.refineArgs(level, args...)
		logger.Entry.Logf(level, format, args...)
	} else {
		logger.Entry.Logf(level, format, args...)
	}
}

func (logger *Logger) refineArgs(level logrus.Level, args ...interface{}) []interface{} {
	for i, arg := range args {
		switch function := arg.(type) {
		case StringGetter:
			args[i] = function()
		}
	}
	return args
}

func getSimplePackageName(pkg interface{}) string {
	return reflect.TypeOf(pkg).PkgPath()
}

func getLogger(module string, levelValue string) *Logger {
	// TODO use Formatter for better log message
	level := getLogLevel(levelValue)
	logger := logrus.New()
	logger.SetLevel(level)
	var logEntry *logrus.Entry
	if len(module) > 0 {
		logEntry = logrus.NewEntry(logger).WithField("module", module)
	} else {
		logEntry = logrus.NewEntry(logger)
	}
	return &Logger{*logEntry}
}

func getLogLevel(levelValue string) logrus.Level {
	level, err := logrus.ParseLevel(levelValue)
	if err != nil {
		level, _ = logrus.ParseLevel(DefaultLogLevel)
	}
	return level
}

func getConfig(module string) Config {
	return config // TODO return log config by module name, fallback to default
}
