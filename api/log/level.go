package log

import "strings"

type Level int8

const LevelKey = "level"

const (
	LevelDebug Level = iota - 1
	LevelInfo
)

func (l Level) String() string {
	switch l {
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	default:
		return ""
	}
}

func ParseLevel(s string) Level {
	switch strings.ToUpper(s) {
	case "DEBUG":
		return LevelDebug
	case "INFO":
		return LevelInfo

	}

	return LevelInfo

}
