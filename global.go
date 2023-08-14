package ilog

import (
	"log"
	"os"
	"sync"
)

//nolint:gochecknoglobals
var (
	globalLogger   Logger = NewBuilder(DebugLevel, os.Stdout).Build() //nolint:revive
	globalLoggerMu sync.RWMutex
)

func Global() Logger {
	globalLoggerMu.RLock()
	defer globalLoggerMu.RUnlock()
	return globalLogger
}

func SetGlobal(logger Logger) (rollback func()) {
	globalLoggerMu.Lock()
	defer globalLoggerMu.Unlock()
	backup := globalLogger

	globalLogger = logger

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
