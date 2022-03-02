package logger

import (
	"fmt"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	zsLogger   *zap.SugaredLogger
	rcvCEnable bool
	rcvC       chan string
}

func (l *Logger) OutputC() <-chan string {
	return l.rcvC
}

var implLogger = &Logger{
	zsLogger:   zap.NewNop().Sugar(),
	rcvCEnable: false,
	rcvC:       make(chan string, 5), // this must be buffered or the hook wii be blocked by channel send
}

func Get() *Logger {
	return implLogger
}

func Set(cfg *LogConfig, enableReceiveC bool) error {
	if cfg == nil {
		return ErrNilConfig
	}

	var (
		zapCfg = zap.NewDevelopmentConfig()
	)
	// log level
	switch strings.ToUpper(cfg.Level) {
	case "FATAL":
		zapCfg.Level.SetLevel(zapcore.FatalLevel)
	case "PANIC":
		zapCfg.Level.SetLevel(zapcore.PanicLevel)
	case "DPANIC":
		zapCfg.Level.SetLevel(zapcore.DPanicLevel)
	case "ERROR":
		zapCfg.Level.SetLevel(zapcore.ErrorLevel)
	case "WARN":
		zapCfg.Level.SetLevel(zapcore.WarnLevel)
	case "INFO":
		zapCfg.Level.SetLevel(zapcore.InfoLevel)
	case "DEBUG":
		zapCfg.Level.SetLevel(zapcore.DebugLevel)
	default:
		zapCfg.Level.SetLevel(zapcore.InfoLevel)
	}

	// encoder
	zapCfg.EncoderConfig = zapcore.EncoderConfig{
		MessageKey: "message",
		// FunctionKey: "function",

		NameKey:    "name",
		EncodeName: zapcore.FullNameEncoder,

		LevelKey:    "level",
		EncodeLevel: zapcore.CapitalLevelEncoder,

		TimeKey:    "time",
		EncodeTime: zapcore.ISO8601TimeEncoder,

		CallerKey:    "caller",
		EncodeCaller: zapcore.ShortCallerEncoder,
	}

	// log color
	if cfg.Color {
		zapCfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	// build
	zLogger, err := zapCfg.Build(zap.AddCallerSkip(1), zap.AddStacktrace(zap.DPanicLevel))
	if err != nil {
		return err
	}

	if enableReceiveC {
		implLogger.rcvCEnable = enableReceiveC
		zLogger = zLogger.WithOptions(zap.Hooks(func(entry zapcore.Entry) error {
			if enableReceiveC {
				implLogger.rcvC <- fmt.Sprintf("[%s]\t%s\t%s", entry.Level.CapitalString(), entry.Caller.Function, entry.Message)
			}
			return nil
		}))
	}
	implLogger.zsLogger = zLogger.Sugar()

	zap.ReplaceGlobals(zLogger.Named(cfg.LoggerName))

	return nil
}

func (l Logger) Fatal(message string) {
	l.zsLogger.Fatal(message)
}

func (l Logger) Fatalf(template string, args ...interface{}) {
	l.zsLogger.Fatalf(template, args...)
}

func (l Logger) Panic(message string) {
	l.zsLogger.Panic(message)
}

func (l Logger) Panicf(template string, args ...interface{}) {
	l.zsLogger.Panicf(template, args...)
}

func (l Logger) DPanic(message string) {
	l.zsLogger.DPanic(message)
}

func (l Logger) DPanicf(template string, args ...interface{}) {
	l.zsLogger.DPanicf(template, args...)
}

func (l Logger) Error(message string) {
	l.zsLogger.Error(message)
}

func (l Logger) Errorf(template string, args ...interface{}) {
	l.zsLogger.Errorf(template, args...)
}

func (l Logger) Warn(message string) {
	l.zsLogger.Warn(message)
}

func (l Logger) Warnf(template string, args ...interface{}) {
	l.zsLogger.Warnf(template, args...)
}

func (l Logger) Info(message string) {
	l.zsLogger.Info(message)
}

func (l Logger) Infof(template string, args ...interface{}) {
	l.zsLogger.Infof(template, args...)
}

func (l Logger) Debug(message string) {
	l.zsLogger.Debug(message)
}

func (l Logger) Debugf(template string, args ...interface{}) {
	l.zsLogger.Debugf(template, args...)
}
