package library

import (
	"context"
	"fmt"
	rotateLog "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"os"
	"strings"
	"talkcheap.xiaoeknow.com/xiaoetong/eframe/contract"
	"time"
)

type Log struct {
	*zap.Logger
}

func (log *Log) HttpRequestLog(record *contract.XiaoeHttpRequestRecord) {
	log.Info(record.Msg,
		zap.String("app_id", record.AppId),
		zap.String("sw8", record.Sw8),
		zap.String("sw8_correlation", record.Sw8Correlation),
		zap.String("xe_tag", record.XeTag),
		zap.String("trace_id", record.TraceId),
		zap.String("uid", record.Uid),
		zap.Int("http_status", record.HttpStatus),
		zap.String("target_url", record.TargetUrl),
		zap.String("method", record.Method),
		zap.String("client_ip", record.ClientIp),
		zap.String("server_ip", record.ServerIp),
		zap.String("user_agent", record.UserAgent),
		zap.String("begin_time", record.BeginTime),
		zap.String("end_time", record.EndTime),
		zap.Int("cost_time", record.CostTime),
		zap.String("params", record.Params),
		zap.String("response", record.Response),
		zap.String("header", record.Header),
	)
}

// eframe-demo有用到
func GetLogger(i interface{}) *Log {
	switch v := i.(type) {
	case *LoggerConfig:
		return NewLogger(v)
	case *LumberLoggerConfig:
		return NewLumberLogger(v)
	case *RotateLoggerConfig:
		return NewRotateLogger(v)
	default:
		panic("unsupported logger conf")
	}
}

func (log *Log) getCtxInfo(ctx context.Context) []zap.Field {
	filed := make([]zap.Field, 0)
	if xeCtx, ok := ctx.Value(contract.XeCtx).(map[string]string); ok {
		filed = append(filed, zap.String("trace_id", xeCtx[contract.TraceId]))
	}
	return filed
}

func (log *Log) InfoCtx(ctx context.Context, msg string, fields ...zap.Field) {
	ctxField := append(log.getCtxInfo(ctx), fields...)
	log.Logger.WithOptions(zap.AddCallerSkip(1)).Info(msg, ctxField...)
}

func (log *Log) WarningCtx(ctx context.Context, msg string, fields ...zap.Field) {
	ctxField := append(log.getCtxInfo(ctx), fields...)
	log.Logger.WithOptions(zap.AddCallerSkip(1)).Warn(msg, ctxField...)
}

func (log *Log) ErrorCtx(ctx context.Context, msg string, fields ...zap.Field) {
	ctxField := append(log.getCtxInfo(ctx), fields...)
	log.Logger.WithOptions(zap.AddCallerSkip(1)).Error(msg, ctxField...)
}

func NewLogger(config *LoggerConfig) *Log {
	fileName := fmt.Sprintf("%s/%s.log", config.Path, config.Name)
	hook := &lumberjack.Logger{
		Filename:   fileName,
		MaxSize:    config.MaxFileSize,
		MaxAge:     config.MaxAgeDay,
		MaxBackups: config.MaxBackups,
		LocalTime:  true,
		Compress:   config.CompressFile,
	}
	baseConf := &BaseLoggerConfig{
		Receiver:       config.Receiver,
		Level:          config.Level,
		Name:           config.Name,
		Path:           config.Path,
		MaxAgeDay:      config.MaxAgeDay,
		MaxFileSize:    config.MaxFileSize,
		MaxBackups:     config.MaxBackups,
		PrintInConsole: config.PrintInConsole,
		DebugMode:      config.DebugMode,
	}
	return getZapLogger(baseConf, hook)
}

func NewLumberLogger(config *LumberLoggerConfig) *Log {
	fileName := fmt.Sprintf("%s/%s.log", config.Path, config.Name)
	hook := &lumberjack.Logger{
		Filename:   fileName,
		MaxSize:    config.MaxFileSize,
		MaxAge:     config.MaxAgeDay,
		MaxBackups: config.MaxBackups,
		LocalTime:  true,
		Compress:   config.CompressFile,
	}
	return getZapLogger(&config.BaseLoggerConfig, hook)
}

func NewRotateLogger(config *RotateLoggerConfig) *Log {
	var MbSize = 1024 * 1024
	fileName := fmt.Sprintf("%s/%s.log", config.Path, config.Name)
	options := []rotateLog.Option{
		rotateLog.WithLinkName(fileName),
		rotateLog.WithRotationSize(int64(config.MaxFileSize * MbSize)),
		rotateLog.WithRotationTime(time.Duration(config.RotationTime) * time.Hour),
	}
	if config.UseMaxBackups {
		options = append(options, rotateLog.WithRotationCount(uint(config.MaxBackups)))
	} else {
		options = append(options, rotateLog.WithMaxAge(time.Duration(config.MaxAgeDay*24)*time.Hour))
	}

	hook, err := rotateLog.New(
		fileName+"_%Y%m%d_%H",
		options...,
	)
	if err != nil {
		panic(err)
	}
	return getZapLogger(&config.BaseLoggerConfig, hook)
}

func getZapLogger(config *BaseLoggerConfig, hook io.Writer) *Log {
	encoderConf := zapcore.EncoderConfig{
		MessageKey:       "msg",
		LevelKey:         "level",
		TimeKey:          "log_at",
		NameKey:          config.Name,
		CallerKey:        "caller",
		FunctionKey:      "function",
		StacktraceKey:    "stack",
		SkipLineEnding:   false,
		LineEnding:       zapcore.DefaultLineEnding,
		EncodeLevel:      zapcore.LowercaseLevelEncoder,
		EncodeTime:       zapcore.ISO8601TimeEncoder,
		EncodeDuration:   zapcore.SecondsDurationEncoder,
		EncodeCaller:     zapcore.FullCallerEncoder,
		EncodeName:       zapcore.FullNameEncoder,
		ConsoleSeparator: "",
	}

	atomicLevel := zap.NewAtomicLevel()
	switch strings.ToLower(config.Level) {
	case "debug":
		atomicLevel.SetLevel(zap.DebugLevel)
	case "info":
		atomicLevel.SetLevel(zap.InfoLevel)
	case "warning":
		atomicLevel.SetLevel(zap.WarnLevel)
	case "error":
		atomicLevel.SetLevel(zap.ErrorLevel)
	default:
		atomicLevel.SetLevel(zap.InfoLevel)
	}

	var ws []zapcore.WriteSyncer
	if config.PrintInConsole {
		ws = append(ws, zapcore.AddSync(os.Stdout))
	}
	ws = append(ws, zapcore.AddSync(hook))

	core := zapcore.NewCore(zapcore.NewJSONEncoder(encoderConf),
		zapcore.NewMultiWriteSyncer(ws...), atomicLevel)

	var options []zap.Option

	if config.DebugMode {
		options = append(options, zap.AddCaller())
		options = append(options, zap.Development())
	}

	filed := zap.Fields(zap.String("service_name", config.ServiceName))

	options = append(options, filed)

	return &Log{
		zap.New(core, options...),
	}
}
