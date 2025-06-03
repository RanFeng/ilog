package ilog

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"encoding/json"
	"runtime/debug"

	"github.com/rs/zerolog"
)

const (
	LogLevelNull = iota // 当前未配置loglevel
	LevelTrace
	LevelDebug
	LevelInfo
	LevelNotice
	LevelWarn
	LevelError
	LevelFatal
)

const (
	LogLevelKey  = "K_LOG_LEVEL"
	LogIDKey     = "K_LOG_ID"
	LogSuffixKey = "K_LOG_SUFFIX" // 后缀在context中的key
	EnvKey       = "RUN_ENV"
)

var (
	logger      zerolog.Logger
	globalLevel = LevelDebug
)

func init() {
	// 创建log目录
	logDir := "./run_log/"
	err := os.MkdirAll(logDir, os.ModePerm)
	if err != nil {
		fmt.Println("Mkdir failed, err:", err)
		return
	}
	fileName := logDir + time.Now().Format("2006-01-02") + ".log"
	logFile, _ := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)

	zerolog.TimeFieldFormat = "2006-01-02 15:04:05.000000"
	consoleWriter := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: "2006-01-02 15:04:05.000000",
	}
	multi := zerolog.MultiLevelWriter(consoleWriter, logFile)
	if env := os.Getenv(EnvKey); env == "prod" {
		// 如果环境变量是prod环境，则只写log文件
		multi = zerolog.MultiLevelWriter(logFile)
	}
	logger = zerolog.New(multi).With().Timestamp().Logger()
}

/* getLevelFromCtx
 * @Description: 从context中获取logLevel
 * @param ctx
 * @return int
 */
func getLevelFromCtx(ctx context.Context) int {
	var ok bool
	var level int
	valObj := ctx.Value(LogLevelKey)
	level, ok = valObj.(int)
	if !ok {
		return globalLevel
	}
	return level
}

/*SetCtxLogLevel
 * @Description: 设置本次请求的context日志级别
 * @param ctx
 * @param int
 */
func SetCtxLogLevel(ctx context.Context, level int) context.Context {
	return context.WithValue(ctx, LogLevelKey, level)
}

/*SetGlobalLogLevel
 * @Description: 设置本服务的context日志级别
 * @param ctx
 * @param int
 */
func SetGlobalLogLevel(level int) {
	globalLevel = level
}

/*CtxLevel
 * @Description: 获取日志级别，优先级：context中携带的配置>全局初始化配置
 * @param ctx
 * @return int 日志级别
 */
func CtxLevel(ctx context.Context) int {
	ctxLevel := getLevelFromCtx(ctx) // 看context配置
	if ctxLevel != LogLevelNull {
		return ctxLevel
	}
	return globalLevel
}

// EventError 记录错误事件日志
func EventError(ctx context.Context, err error, event string, payload ...interface{}) {
	errPayload := []interface{}{"err", "nil"}
	if err != nil {
		errPayload[1] = err.Error()
	}
	payload = append(errPayload, payload...)
	doLog(ctx, LevelError, event, payload...)
}

// EventFatal 记录严重事故事件日志
func EventFatal(ctx context.Context, event string, payload ...interface{}) {
	doLog(ctx, LevelFatal, event, payload...)
}

// EventInfo 记录Info事件日志
func EventInfo(ctx context.Context, event string, payload ...interface{}) {
	doLog(ctx, LevelInfo, event, payload...)
}

// EventDebug 记录调试事件日志
func EventDebug(ctx context.Context, event string, payload ...interface{}) {
	doLog(ctx, LevelDebug, event, payload...)
}

// EventWarn  记录Warn事件日志
func EventWarn(ctx context.Context, event string, payload ...interface{}) {
	doLog(ctx, LevelWarn, event, payload...)
}

// doLog 完成日志内容格式化及写入
// 日志内容包括四个部分：
// 1. 常规日志信息信息，例如时间、代码行等
// 2. 本条日志的event、error、payload等信息
// 3. 日志上下文中存储的公共信息，例如AppID等
// 4. 日志上下文中存储的定制化信息
// 以上数据按顺序拼接，成为日志的最终内容
func doLog(ctx context.Context, logLevel int, event string, kvList ...interface{}) {
	if logLevel < CtxLevel(ctx) {
		return
	}
	defer func() {
		if err := recover(); err != nil {
			logID, _ := ctx.Value(LogIDKey).(string)
			ctx = context.WithValue(context.Background(), LogIDKey, logID)
			EventFatal(ctx, "ilog_panic", "stack", debug.Stack())
		}
	}()

	var logObj *zerolog.Event
	switch logLevel {
	case LevelFatal:
		logObj = logger.Fatal()
	case LevelError:
		logObj = logger.Error()
	case LevelWarn:
		logObj = logger.Warn()
	case LevelInfo:
		logObj = logger.Info()
	case LevelDebug:
		logObj = logger.Debug()
	default:
		logObj = logger.Error()
	}
	// 打印logid
	logID, _ := ctx.Value(LogIDKey).(string)
	logObj.Str("log_id", logID)
	logObj.Str("event", event)
	payload := zerolog.Dict()
	if len(kvList)%2 == 1 {
		kvList = append(kvList, "")
	}
	for i := 0; i+1 < len(kvList); i += 2 {
		payload.Str(InterfaceToString(kvList[i]), InterfaceToString(kvList[i+1]))
	}
	logObj.Dict("payload", payload)
	valObj := ctx.Value(LogSuffixKey)
	logObj.Any("suffix", valObj)
	logObj.Caller(2).Send()
}

type LogStringer interface {
	LogString() string
}

/*
InterfaceToString
  - @Description: 转化各种类型到string，
    注意！如果想转化float到string，一般有两种方式：
    只关注float的整数部分，
    InterfaceToString(3.124567889) = 3
    InterfaceToString(23412353.12412) = 23412353
    其中val是想要转化的float
  - @param value
  - @return string
*/
func InterfaceToString(value interface{}) string {
	if value == nil {
		return ""
	}
	switch v := value.(type) {
	case string:
		return v
	case fmt.GoStringer:
		return v.GoString()
	case LogStringer:
		return v.LogString()
	case bool:
		if v {
			return "true"
		} else {
			return "false"
		}
	case error:
		return v.Error()

	case int:
		return strconv.FormatInt(int64(v), 10)
	case int8:
		return strconv.FormatInt(int64(v), 10)
	case int16:
		return strconv.FormatInt(int64(v), 10)
	case int32:
		return strconv.FormatInt(int64(v), 10)
	case int64:
		return strconv.FormatInt(v, 10)

	case uint:
		return strconv.FormatUint(uint64(v), 10)
	case uint8:
		return strconv.FormatUint(uint64(v), 10)
	case uint16:
		return strconv.FormatUint(uint64(v), 10)
	case uint32:
		return strconv.FormatUint(uint64(v), 10)
	case uint64:
		return strconv.FormatUint(v, 10)

	case float32:
		return strconv.FormatInt(int64(v), 10)
	case float64:
		return strconv.FormatInt(int64(v), 10)
	default:
		bytes, err := json.Marshal(v)
		if err != nil {
			return ""
		}
		return string(bytes)
	}
}
