package ulog

import (
	"fmt"
)

// 颜色优先级从高到低:
// 1. 开启颜色(EnableColor)
// 2. 指定颜色(WithColorName, WithColor)
// 3. 日志文本匹配规则(ParseColorRule)
// 4. 日志级别对应颜色(GetColorDefineByLevel)
type TextFormatter struct {
	EnableColor     bool // 命令行着色
	TimestampFormat string

	rule *ColorRuleSet
}

// 从颜色规则文本读取规则
func (self *TextFormatter) ParseColorRule(ruleText string) error {
	self.rule = NewColorRuleSet()
	err := self.rule.Parse(ruleText)
	if err != nil {
		return fmt.Errorf("parse color rule failed: %w", err)
	}

	return nil
}

func (self *TextFormatter) matchText(text string) *ColorDefine {

	if self.rule == nil {
		return nil
	}

	return self.rule.MatchText(text)
}

// 取得颜色前缀
func (self *TextFormatter) GetPrefix(entry *Entry) *ColorDefine {

	var cdef *ColorDefine

	if self.EnableColor {

		cdef = entry.ColorDef

		if cdef == nil {
			cdef = self.matchText(entry.Message)

			if cdef == nil {
				cdef = GetColorDefineByLevel(entry.Level)
			}
		}
	} else {
		cdef = WhiteColorDef
	}

	return cdef
}

// 取得颜色后缀
func (self *TextFormatter) GetSuffix() string {
	if self.EnableColor {
		return consoleColorSuffix
	}
	return ""
}

// 取得时间
func (self *TextFormatter) GetTime(entry *Entry) string {
	var timeFormat string
	if self.TimestampFormat != "" {
		timeFormat = self.TimestampFormat
	} else {
		timeFormat = TextTimeFormat
	}
	return entry.Time.Format(timeFormat)
}

// 取得调用者
func (self *TextFormatter) GetCaller() string {
	return ""
}

func (self *TextFormatter) Format(entry *Entry) ([]byte, error) {

	b := entry.Buffer

	fmt.Fprintf(b, "%s%s[%s]%s %s%s",
		self.GetPrefix(entry),
		entry.Level.String(),
		self.GetTime(entry),
		self.GetCaller(),
		entry.Message,
		self.GetSuffix(),
	)

	// TODO
	// 根据配置添加kv

	b.WriteByte('\n')
	return b.Bytes(), nil
}
