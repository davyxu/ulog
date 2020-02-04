package ulog

type Level uint32

const (
	DebugLevel Level = iota
	InfoLevel
	WarnLevel
	ErrorLevel
)

func (self Level) String() string {
	return levelString[self]
}

var levelString = [...]string{
	"DEBU",
	"INFO",
	"WARN",
	"ERRO",
}
