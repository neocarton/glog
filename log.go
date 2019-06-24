package glog

import (
	"encoding/json"
	"io/ioutil"
	"reflect"
	"strings"
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
	StringGetter func() string

	// Config interface
	Config interface {
		Package() string
		Level() string
	}

	// DefaultConfig default config
	DefaultConfig struct {
		Config
		pkgName string
		level   string
	}

	// Logger logger
	Logger struct {
		logrus.Entry
	}
)

var (
	config Config = &DefaultConfig{level: DefaultLogLevel}
)

// Initialize logger
func Initialize(cfg Config) {
	config = cfg
	level := getLogLevel(config.Level())
	logrus.SetLevel(level)
}

// AsJSON Convert object to JSON when needed
func AsJSON(object interface{}) StringGetter {
	return func() string {
		if object == nil {
			return ""
		}
		data, _ := json.MarshalIndent(object, "", "  ")
		return string(data)
	}
}

// AsISOTime Convert object to JSON when needed
func AsISOTime(t time.Time) StringGetter {
	return func() string {
		return t.Format(time.RFC3339)
	}
}

// Package return package name
func (cfg *DefaultConfig) Package() string {
	return cfg.pkgName
}

// SetPackage set package name
func (cfg *DefaultConfig) SetPackage(pkgName string) {
	cfg.pkgName = pkgName
}

// Level return log level
func (cfg *DefaultConfig) Level() string {
	return cfg.pkgName
}

// SetLevel return log level
func (cfg *DefaultConfig) SetLevel(level string) {
	cfg.level = level
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
		format += ": %s" // TODO print log trace
		args = append(args, AsJSON(err))
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
		case StringGetter: // TODO Support more getter
			args[i] = function()
		}
	}
	return args
}

func getSimplePackageName(pkg interface{}) string {
	pkgName := reflect.TypeOf(pkg).PkgPath()
	pkgPrefix := config.Package()
	if len(pkgPrefix) > 0 && strings.HasPrefix(pkgName, pkgPrefix) {
		pkgName = pkgName[len(pkgPrefix):]
	}
	return pkgName
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
