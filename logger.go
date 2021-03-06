package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var defaultLogger *zap.Logger

func getLogger() *zap.Logger {
	if defaultLogger == nil {
		err := InitLogger(false, "", "debug") // 默认输出到控台
		if err != nil {
			log.Fatal(err)
		}
	}

	return defaultLogger.WithOptions(zap.AddCallerSkip(1))
}

// InitLogger 初始化日志
//	isSave 是否输出到文件，true: 是，false:输出到控台
//	filename 保存日志路径，例如："out.log"
//	level  输出日志级别 DEBUG, INFO, WARN, ERROR
//	encodingType  输出格式 json:显示数据格式为json，console:显示数据格式为console(默认)
// 		以console数据格式输出到控台，eg: InitLogger(false, "", "debug")
// 		以json数据格式输出到控台，eg: InitLogger(false, "", "debug", "json")
// 		以json数据格式输出到文件，eg: InitLogger(true, "out.log", "debug")
func InitLogger(isSave bool, filename string, level string, encodingType ...string) error {
	log.SetFlags(log.Lmicroseconds | log.Lshortfile | log.LstdFlags) // log包显示设置

	// 保存日志路径
	if isSave && filename == "" {
		filename = "out.log" // 默认
	}

	// 日志输出等级
	levelName := ""
	switch strings.ToUpper(level) {
	case "DEBUG":
		levelName = "DEBUG"
	case "INFO":
		levelName = "INFO"
	case "WARN":
		levelName = "WARN"
	case "ERROR":
		levelName = "ERROR"
	default:
		levelName = "DEBUG" // 默认
	}

	var encoding string
	var js string
	if isSave { // 日志保存到文件
		encoding = "json" // 当日志输出到文件时，只有json格式

		js = fmt.Sprintf(`{
      		"level": "%s",
      		"encoding": "%s",
      		"outputPaths": ["%s"],
      		"errorOutputPaths": ["%s"]
      	}`, levelName, encoding, filename, filename)
	} else { // 在控台输出日志
		if len(encodingType) > 0 && encodingType[0] == "json" { // 控台模式下可以输出json格式，也可以输出console模式
			encoding = "json"
		} else {
			encoding = "console"
		}

		js = fmt.Sprintf(`{
      		"level": "%s",
            "encoding": "%s",
      		"outputPaths": ["stdout"],
            "errorOutputPaths": ["stdout"]
		}`, levelName, encoding)
	}

	var config zap.Config
	err := json.Unmarshal([]byte(js), &config)
	if err != nil {
		return err
	}

	config.EncoderConfig = zap.NewProductionEncoderConfig()

	config.EncoderConfig.EncodeTime = timeFormatter // 默认时间格式
	if isSave {
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	}

	defaultLogger, err = config.Build()
	if err != nil {
		return err
	}

	// 打印log配置结果
	if isSave {
		Infof("initialize logger finish, base config is isSave=%t, filename=%s, level=%s, encoding=%s", isSave, filename, level, encoding)
	} else {
		Infof("initialize logger finish, base config is isSave=%t, level=%s, encoding=%s", isSave, level, encoding)
	}

	return nil
}

func timeFormatter(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}

// Ctx logs trace info
// X-B3-TraceId：一条请求链路（Trace）的唯一标识，必须值
// X-B3-SpanId：一个工作单元（Span）的唯一标识，必须值
// X-B3-ParentSpanId:：标识当前工作单元所属的上一个工作单元，Root Span（请求链路的第一个工作单元）的该值为空
// X-B3-Sampled：是否被抽样输出的标志，1表示需要被输出，0表示不需要被输出
// X-Span-Name：工作单元的名称
func Ctx(ctx context.Context) *zap.Logger {
	fieldsMap := make(map[string]interface{})
	keys := []string{"X-B3-TraceId", "X-B3-SpanId", "X-B3-ParentSpanId", "X-Span-Name"}

	if ctx != nil {
		for _, key := range keys {
			if v := ctx.Value(key); v != nil {
				fieldsMap[key] = v
			}
		}
	}

	if len(fieldsMap) > 0 {
		return getLogger().With(Any("context", fieldsMap))
	}

	return getLogger()
}

// ----------------------------------重新封装zap的log----------------------------------------

// Debug debug级别信息
func Debug(msg string, fields ...Field) {
	getLogger().Debug(msg, fields...)
}

// Info info级别信息
func Info(msg string, fields ...Field) {
	getLogger().Info(msg, fields...)
}

// Warn warn级别信息
func Warn(msg string, fields ...Field) {
	getLogger().Warn(msg, fields...)
}

// Error error级别信息
func Error(msg string, fields ...Field) {
	getLogger().Error(msg, fields...)
}

// Panic panic级别信息
func Panic(msg string, fields ...Field) {
	getLogger().Panic(msg, fields...)
}

// Fatal fatal级别信息
func Fatal(msg string, fields ...Field) {
	getLogger().Fatal(msg, fields...)
}

// Debugf 带格式化debug级别信息
func Debugf(format string, a ...interface{}) {
	getLogger().Debug(fmt.Sprintf(format, a...))
}

// Infof 带格式化info级别信息
func Infof(format string, a ...interface{}) {
	getLogger().Info(fmt.Sprintf(format, a...))
}

// Warnf 带格式化warn级别信息
func Warnf(format string, a ...interface{}) {
	getLogger().Warn(fmt.Sprintf(format, a...))
}

// Errorf 带格式化error级别信息
func Errorf(format string, a ...interface{}) {
	getLogger().Error(fmt.Sprintf(format, a...))
}

// Fatalf 带格式化fatal级别信息
func Fatalf(format string, a ...interface{}) {
	getLogger().Fatal(fmt.Sprintf(format, a...))
}

// WithFields 携带字段信息
func WithFields(fields ...Field) *zap.Logger {
	return getLogger().With(fields...)
}

// ----------------------- 重新封装类型 ---------------------------

// ZapLogger logger类型
type ZapLogger = zap.Logger

// Field 字段类型
type Field = zapcore.Field

// Int int类型
func Int(key string, val int) Field {
	return zap.Int(key, val)
}

// Int64 int64类型
func Int64(key string, val int64) Field {
	return zap.Int64(key, val)
}

// Uint uint类型
func Uint(key string, val uint) Field {
	return zap.Uint(key, val)
}

// Uint64 uint64类型
func Uint64(key string, val uint64) Field {
	return zap.Uint64(key, val)
}

// Uintptr uintptr类型
func Uintptr(key string, val uintptr) Field {
	return zap.Uintptr(key, val)
}

// Float64 float64类型
func Float64(key string, val float64) Field {
	return zap.Float64(key, val)
}

// Bool bool类型
func Bool(key string, val bool) Field {
	return zap.Bool(key, val)
}

// String string类型
func String(key string, val string) Field {
	return zap.String(key, val)
}

// Stringer stringer类型
func Stringer(key string, val fmt.Stringer) Field {
	return zap.Stringer(key, val)
}

// Time time.Time类型
func Time(key string, val time.Time) Field {
	return zap.Time(key, val)
}

// Duration time.Duration类型
func Duration(key string, val time.Duration) Field {
	return zap.Duration(key, val)
}

// Err err类型
func Err(err error) Field {
	return zap.Error(err)
}

// Any 任意类型，如果是对象、slice、map等复合类型，使用Any
func Any(key string, val interface{}) Field {
	return zap.Any(key, val)
}

// GetLogger 获取defaultLogger，设置caller值才能正确的显示对应的代码行数
func GetLogger(skip int) *zap.Logger {
	if defaultLogger == nil {
		err := InitLogger(false, "", "debug") // 默认输出到控台
		if err != nil {
			log.Fatal(err)
		}
	}

	return defaultLogger.WithOptions(zap.AddCallerSkip(skip))
}