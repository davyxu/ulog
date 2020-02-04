package ulog

type Level uint32

const (
	DebugLevel Level = iota
	InfoLevel
	WarnLevel
	ErrorLevel
)

func (self Level) String() string {
	switch self {
	case DebugLevel:
		return "DEBU"
	case InfoLevel:
		return "INFO"
	case WarnLevel:
		return "WARN"
	case ErrorLevel:
		return "ERRO"
	}

	return "UNKW"
}

var levelBytes = [...][]byte{
	[]byte("DEBU"),
	[]byte("INFO"),
	[]byte("WARN"),
	[]byte("ERRO"),
}
