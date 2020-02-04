package ulog

import (
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestTextFormatter(t *testing.T) {

	Global().SetFormatter(&TextFormatter{
		EnableColor: true,
	})

	// 基本用法
	Global().SetLevel(DebugLevel)
	Debugln("debug")
	Infof("info")
	Warnf("warning %d", 123)
	Errorf("error %d", 567)

	// 全局颜色输出
	WithColorName("purple").Infoln("WithColorName ", Purple.String())
	WithColor(DarkGreen).Errorf("WithColor %s", DarkGreen.String())

	// 独立日志实例
	l := New()

	textFormat := &TextFormatter{
		EnableColor: true,
	}

	// 从文本解析
	err := textFormat.ParseColorRule(`
{
	"Rule":[
		{"Text":"panic:","Color":"Red"},
		{"Text":"[DB]","Color":"Green"},
		{"Text":"#http.listen","Color":"Blue"},
		{"Text":"#http.recv","Color":"Blue"},
		{"Text":"#http.send","Color":"Purple"}
	]
}
`)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	l.SetFormatter(textFormat)

	// 从文本确定颜色
	l.Infoln("panic: must be red")
	l.Infoln("[DB] write db")
	l.Infoln("#http.recv come data")
}

func TestJsonFormatter(t *testing.T) {
	Global().SetFormatter(&JSONFormatter{})
	Global().SetLevel(DebugLevel)

	// 单行kv
	WithField("key", "value").Infof("noraml json")

	Global().SetFormatter(&JSONFormatter{
		PrettyPrint: true,
	})

	// 多行kv
	Global().WithFields(Fields{
		"name": "monk",
		"age":  80,
	}).Errorf("error json with pretty print")
}

// 定义纯文本输出
type testFormatter struct {
	TextFormatter
}

func (self *testFormatter) Format(entry *Entry) ([]byte, error) {
	b := entry.Buffer

	b.WriteString(entry.Message)

	b.WriteByte('\n')
	return b.Bytes(), nil
}

func TestRollingFile(t *testing.T) {
	const (
		fileName        = "tt.log"
		maxFileSize     = 1000
		eachTimeWrite   = 100
		totalTimesWrite = 100
		numberFile      = eachTimeWrite * totalTimesWrite / maxFileSize
	)

	Global().SetFormatter(&testFormatter{})

	asyncWriter := NewAsyncOutput(NewRollingOutput(fileName, maxFileSize))
	Global().SetOutput(asyncWriter)

	for i := 0; i < totalTimesWrite; i++ {
		Infoln(strings.Repeat(strconv.Itoa(i), eachTimeWrite))
	}

	// 异步写入时, 在程序结束前, 需要保证完全写入
	asyncWriter.Flush(time.Second)
}

// 没有Entry分配
func BenchmarkNoEntryAlloc(b *testing.B) {

	for i := 0; i < 10; i++ {
		Infoln(i)
	}

}

// 独立Entry不进内存池, 有多少分配多少
func BenchmarkEntryEach(b *testing.B) {

	for i := 0; i < 10; i++ {
		NewEntry(Global()).Infoln(i)
	}
}