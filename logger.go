package toolkit

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/valyala/bytebufferpool"
)

// ============ 日志 ============ //

type level uint8

const (
	LevelDebug level = iota + 1
	LevelInfo
	LevelWarn
	LevelError
)

type Logger struct {
	w     io.Writer
	level level
}

func NewLogger(filename string, level level) *Logger {
	l := &Logger{
		w:     os.Stdout,
		level: level,
	}

	if filename != "" {
		var err error
		l.w, err = FileOpen(filename)
		if err != nil {
			panic("open log file error: " + err.Error())
		}
	}

	return l
}

func (l *Logger) Debug(args ...any) {
	if l.level >= LevelDebug && len(args) > 0 {
		write(l, LevelDebug, args...)
	}
}

func (l *Logger) Info(args ...any) {
	if l.level >= LevelInfo && len(args) > 0 {
		write(l, LevelInfo, args...)
	}
}

func (l *Logger) Warn(args ...any) {
	if l.level >= LevelWarn && len(args) > 0 {
		write(l, LevelWarn, args...)
	}
}

func (l *Logger) Error(args ...any) {
	if l.level >= LevelError && len(args) > 0 {
		write(l, LevelError, args...)
	}
}

var (
	logBufferPool = &bytebufferpool.Pool{} // buffer池
	levelAlias    = [...]string{"", "[DEBUG]", "[INFO]", "[WARN]", "[ERROR]"}
)

func write(l *Logger, level level, args ...any) {
	buf := logBufferPool.Get()
	defer logBufferPool.Put(buf)
	buf.WriteString(time.Now().Format(time.DateTime))
	buf.WriteString(" " + levelAlias[level] + " ")

	for _, arg := range args {
		switch v := arg.(type) {
		case []byte:
			buf.Write(v)
		case string:
			buf.WriteString(v)
		case uint8, uint16, uint32, uint64, int8, int16, int32, int64, float32, float64, bool:
			buf.WriteString(fmt.Sprintf("%v", v))
		case error:
			buf.WriteString(v.Error())
		default:
			b, err := json.Marshal(arg)
			if err != nil {
				buf.WriteString(fmt.Sprintf("%v", v))
			} else {
				buf.Write(b)
			}
		}
		buf.WriteByte(' ')
	}
	buf.WriteByte('\n')

	l.w.Write(buf.Bytes())
}

// ============ 默认日志 ============ //

var (
	defaultLogger *Logger = NewLogger("", LevelDebug)
)

func SetDefaultLogger(l *Logger) {
	defaultLogger = l
}

func Debug(args ...any) {
	defaultLogger.Debug(args...)
}

func Info(args ...any) {
	defaultLogger.Info(args...)
}

func Warn(args ...any) {
	defaultLogger.Warn(args...)
}

func Error(args ...any) {
	defaultLogger.Error(args...)
}
