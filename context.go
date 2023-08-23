package ilog

import "context"

type contextKey string

const contextKeyLogger contextKey = "Logger"

func FromContext(ctx context.Context) (logger Logger) {
	if ctx == nil {
		Global().Copy().AddCallerSkip(1).Errorf("ilog: nil context")
		return Global().Copy()
	}

	v := ctx.Value(contextKeyLogger)
	l, ok := v.(Logger)

	if !ok {
		Global().Copy().AddCallerSkip(1).Errorf("ilog: type assertion failed: expected=ilog.Logger, actual=%T, value=%#v", v, v)
		return Global().Copy()
	}

	return l.Copy()
}

func WithContext(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, contextKeyLogger, logger)
}
