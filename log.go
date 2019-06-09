package glog

import (
	"io/ioutil"
	"reflect"
	"strings"

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
	// Config structure
	Config struct {
		Package string
		Level   string `json:"level"`
	}

	// Logger logger
	Logger struct {
		logrus.Entry
	}
)

var (
	config = &Config{Level: DefaultLogLevel}
)

// Initialize logger
func Initialize(cfg *Config) {
	config = cfg
	level := getLogLevel(config.Level)
	logrus.SetLevel(level)
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
		level = cfg.Level
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
	return getLogger(module, cfg.Level)
}

// Warnf log error
func (logger *Logger) Warnf(format string, err error, args ...interface{}) {
	// TODO print log trace
	format += ": %+v"
	logger.Entry.Warnf(format, args, err)
}

// Errorf log error
func (logger *Logger) Errorf(format string, err error, args ...interface{}) {
	// TODO print log trace
	format += ": %+v"
	logger.Entry.Errorf(format, args, err)
}

// Fatalf log error
func (logger *Logger) Fatalf(format string, err error, args ...interface{}) {
	// TODO print log trace
	format += ": %+v"
	logger.Entry.Fatalf(format, args, err)
}

func getSimplePackageName(pkg interface{}) string {
	pkgName := reflect.TypeOf(pkg).PkgPath()
	if len(config.Package) > 0 && strings.HasPrefix(pkgName, config.Package) {
		pkgName = pkgName[len(config.Package):]
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

func getConfig(module string) *Config {
	return config // TODO return log config by module name, fallback to default
}
