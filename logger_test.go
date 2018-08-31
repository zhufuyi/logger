package logger

import (
	"context"
	"testing"

	"github.com/iris-contrib/errors"
)

type people struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func init() {
	InitLogger(false, "", "debug") // 以console数据格式输出到控台
	//InitLogger(false, "", "debug", "json") // 以json数据格式输出到控台
	//InitLogger(true, "plugins/out.plugins", "debug") // 以json数据格式输出到文件
}

func TestInitLogger(t *testing.T) {

	logger.Debug("msg", String("key", "this is debug"))
	logger.Info("msg", String("key", "this is info"))
	logger.Warn("msg", String("key", "this is warn"))
	logger.Error("msg", String("key", "this is error"))

	p := &people{"张三", 11}
	ps := []people{{"张三", 11}, {"李四", 12}}
	pMap := map[string]people{"123": *p, "456": *p}
	logger.Info("msg:", Any("object", p))
	logger.Info("msg:", Any("object", ps))
	logger.Info("msg:", Any("object", pMap))

	logger.Error("msg", Err(errors.New("this is error")))
}

func TestPackageLog(t *testing.T) {
	Debug("this is debug")
	Info("this is info")
	Warn("this is warn")
	Error("this is error")

	p := &people{"张三", 11}
	ps := []people{{"张三", 11}, {"李四", 12}}
	pMap := map[string]people{"123": *p, "456": *p}
	Debug("this is debug object", Any("object1", p), Any("object2", ps), Any("object3", pMap))
	Error("err is not equal nil ", Any("object", ps))
	logger.With(Int("hight", 170)).Debug("msg ", Any("object", p))

	ctx := context.WithValue(context.Background(), "X-B3-TraceId", "123456")
	ctx = context.WithValue(ctx, "X-B3-SpanId", "abcdef")
	ctx = context.WithValue(ctx, "X-B3-ParentSpanId", "1a2b3c")
	ctx = context.WithValue(ctx, "X-Span-Name", "logger test")
	Ctx(ctx).Debug("this is debug", Any("object", ps))
}

func BenchmarkString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Info("this is info", String("string", "hello golang"))
	}
}

func BenchmarkInt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Info("this is info", Int("int", 1000))
	}
}

func BenchmarkErr(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Info("this is info", Err(errors.New("error")))
	}
}

func BenchmarkAny(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Info("this is info", Any("object", &people{"张三", 11}))
	}
}
