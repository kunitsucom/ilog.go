package ilog

import (
	"log"
	"os"
	"sync"
)

//nolint:gochecknoglobals
var (
	_globalLogger   Logger = NewBuilder(DebugLevel, os.Stdout).Build() //nolint:revive
	_globalLoggerMu sync.RWMutex
)

func L() Logger {
	_globalLoggerMu.RLock()
	l := _globalLogger
	_globalLoggerMu.RUnlock()
	return l
}

func Global() Logger {
	return L()
}

func SetGlobal(logger Logger) (rollback func()) {
	_globalLoggerMu.Lock()
	backup := _globalLogger
	_globalLogger = logger
	_globalLoggerMu.Unlock()
	return func() {
		SetGlobal(backup)
	}
}

//nolint:gochecknoglobals
var stdLoggerMu sync.Mutex

func SetStdLogger(l Logger) (rollback func()) {
	stdLoggerMu.Lock()
	defer stdLoggerMu.Unlock()

	backupFlags := log.Flags()
	backupPrefix := log.Prefix()
	backupWriter := log.Writer()

	log.SetFlags(0)
	log.SetPrefix("")
	log.SetOutput(l.Copy().AddCallerSkip(2))

	return func() {
		log.SetFlags(backupFlags)
		log.SetPrefix(backupPrefix)
		log.SetOutput(backupWriter)
	}
}
