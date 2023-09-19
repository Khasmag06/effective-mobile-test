package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"strings"
)

type Logger struct {
	logger *zap.SugaredLogger
}

func New(logFilePath, level string) (*Logger, error) {
	var l zapcore.Level

	switch strings.ToLower(level) {
	case "error":
		l = zapcore.ErrorLevel
	case "warn":
		l = zapcore.WarnLevel
	case "info":
		l = zapcore.InfoLevel
	case "debug":
		l = zapcore.DebugLevel
	default:
		l = zapcore.InfoLevel
	}

	//config := zap.NewProductionConfig()
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("Jan 02 15:04:05")
	config.EncoderConfig.CallerKey = zapcore.OmitKey
	config.EncoderConfig.StacktraceKey = zapcore.OmitKey

	config.Level = zap.NewAtomicLevelAt(l)

	if logFilePath != "" {
		config.OutputPaths = []string{logFilePath}
	}

	logger, err := config.Build()
	if err != nil {
		return nil, err
	}
	sugar := logger.Sugar()

	return &Logger{
		logger: sugar,
	}, nil
}

func (l *Logger) Debug(args ...any) {
	l.logger.Debug(args)
}

func (l *Logger) Info(args ...any) {
	l.logger.Info(args)
}

func (l *Logger) Infof(format string, args ...any) {
	l.logger.Infof(format, args)
}

func (l *Logger) Warn(args ...any) {
	l.logger.Warn(args)
}

func (l *Logger) Error(args ...any) {
	l.logger.Error(args)
}

func (l *Logger) Errorf(format string, args ...any) {
	l.logger.Errorf(format, args)
}

func (l *Logger) Fatal(args ...any) {
	l.logger.Fatal(args)
}

func (l *Logger) Fatalf(format string, args ...any) {
	l.logger.Fatalf(format, args)
}

func (l *Logger) Sync() error {
	err := l.logger.Sync()
	if err != nil {
		return err
	}
	return nil
}
