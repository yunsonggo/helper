package loggers

import (
	"github.com/yunsonggo/helper/types"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

func NewLogger(conf *types.Log) (*zap.Logger, error) {
	if conf == nil {
		logConfig := &types.Log{
			Level:        "debug",
			Path:         "logs/access.log",
			MaxSizeMB:    20,
			MaxAgeDay:    30,
			MaxBackupDay: 7,
			Compress:     false,
		}
		conf = logConfig
	}
	encoder := logEncoder()
	writer := logWriter(conf)
	level := new(zapcore.Level)
	if err := level.UnmarshalText([]byte(conf.Level)); err != nil {
		return nil, err
	}
	var core zapcore.Core
	if conf.Level == "debug" {
		consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
		core = zapcore.NewTee(
			zapcore.NewCore(encoder, writer, level),
			zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), zapcore.DebugLevel),
		)
	} else {
		core = zapcore.NewCore(encoder, writer, level)
	}
	caller := zap.AddCaller()
	lg := zap.New(core, caller)
	return lg, nil
}

func logEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	return zapcore.NewJSONEncoder(encoderConfig)
}

func logWriter(conf *types.Log) zapcore.WriteSyncer {
	lj := &lumberjack.Logger{
		Filename:   conf.Path,
		MaxSize:    conf.MaxSizeMB,
		MaxBackups: conf.MaxBackupDay,
		MaxAge:     conf.MaxAgeDay,
		Compress:   conf.Compress,
	}
	return zapcore.AddSync(lj)
}
