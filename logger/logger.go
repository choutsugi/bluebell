package logger

import (
	"bluebell/setting"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

func InitLogger(config *setting.LogConfig) (err error) {

	level, err := getLevel(config.Level)
	if err != nil {
		return
	}

	fileWriter := getFileWriter(config.FileName, config.MaxSize, config.MaxBackups, config.MaxAge)
	stdoutWriter := getStdoutWriter()

	jsonEncoder := getJsonEncoder()
	consoleEncoder := getConsoleEncoder()

	core := zapcore.NewTee(
		zapcore.NewCore(jsonEncoder, fileWriter, *level),
		zapcore.NewCore(consoleEncoder, stdoutWriter, *level),
	)

	logger := zap.New(core, zap.AddCaller())
	zap.ReplaceGlobals(logger)
	return
}

func Sync() {
	zap.L().Sync()
}

func getLevel(levelStr string) (level *zapcore.Level, err error) {
	level = new(zapcore.Level)
	if err = level.UnmarshalText([]byte(levelStr)); err != nil {
		return nil, err
	}
	return
}

func getFileWriter(filename string, maxSize, maxBackup, maxAge int) zapcore.WriteSyncer {
	lumberjackLogger := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    maxSize,
		MaxAge:     maxAge,
		MaxBackups: maxBackup,
	}

	return zapcore.AddSync(lumberjackLogger)
}

func getStdoutWriter() zapcore.WriteSyncer {
	return zapcore.AddSync(zapcore.Lock(os.Stdout))
}

func getJsonEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	return zapcore.NewJSONEncoder(encoderConfig)
}

func getConsoleEncoder() zapcore.Encoder {
	return zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
}
