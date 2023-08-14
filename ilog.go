package ilog

import (
	"errors"
	"io"
	"time"
)

var ErrLogEntryIsNotWritten = errors.New("ilog: log entry not written")

type Level int8

const (
	DebugLevel Level = -8
	InfoLevel  Level = 0
	WarnLevel  Level = 8
	ErrorLevel Level = 16
)

type Logger interface {
	Level() (currentLoggerLevel Level)
	SetLevel(level Level) (logger Logger)
	AddCallerSkip(skip int) (logger Logger)
	Copy() (copiedLogger Logger)

	common
}

type common interface {
	Any(key string, value interface{}) (entry LogEntry)
	Bool(key string, value bool) (entry LogEntry)
	Bytes(key string, value []byte) (entry LogEntry)
	Duration(key string, value time.Duration) (entry LogEntry)
	Err(err error) (entry LogEntry)
	ErrWithKey(key string, err error) (entry LogEntry)
	Float32(key string, value float32) (entry LogEntry)
	Float64(key string, value float64) (entry LogEntry)
	Int(key string, value int) (entry LogEntry)
	Int32(key string, value int32) (entry LogEntry)
	Int64(key string, value int64) (entry LogEntry)
	String(key, value string) (entry LogEntry)
	Time(key string, value time.Time) (entry LogEntry)
	Uint(key string, value uint) (entry LogEntry)
	Uint32(key string, value uint32) (entry LogEntry)
	Uint64(key string, value uint64) (entry LogEntry)

	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Logf(level Level, format string, args ...interface{})
	io.Writer
}

type LogEntry interface {
	common

	Logger() (copiedLogger Logger)
	error // Consider undispatched LogEntry as error so that they can be detected by Go static analysis.
}
