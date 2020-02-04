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

func (self *TextFormatter) Format(entry *Entry) ([]byte, error) {

	b := entry.Buffer

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

	var timeFormat string
	if self.TimestampFormat != "" {
		timeFormat = self.TimestampFormat
	} else {
		timeFormat = TextTimeFormat
	}

	var caller string

	fmt.Fprintf(b, "%s%s[%s]%s %s", cdef.Prefix, levelBytes[entry.Level], entry.Time.Format(timeFormat), caller, entry.Message)

	// TODO
	// 根据配置添加kv

	if self.EnableColor {
		b.Write(consoleColorSuffix)
	}

	b.WriteByte('\n')
	return b.Bytes(), nil
}
