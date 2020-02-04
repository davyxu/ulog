package ulog

import "testing"

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
