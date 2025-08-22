package errs

import (
	"fmt"
	"runtime"
	"strings"
	"time"
)

type ServerError struct {
	Message    string
	Code       string
	Timestamp  time.Time
	StackTrace string
}

func (e *ServerError) Error() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Error: %s", e.Message))
	if e.Code != "" {
		sb.WriteString(fmt.Sprintf(" [Code: %s]", e.Code))
	}
	sb.WriteString(fmt.Sprintf(" [Time: %s]", e.Timestamp.Format(time.RFC3339)))
	if e.StackTrace != "" {
		sb.WriteString(fmt.Sprintf("\nStack Trace:\n%s", e.StackTrace))
	}
	return sb.String()
}

func New(message string, opts ...Option) error {
	cfg := &config{
		includeStack: true,
	}

	for _, opt := range opts {
		opt(cfg)
	}

	var stackTrace string
	if cfg.includeStack {
		stackTrace = captureStackTrace(2)
	}

	return &ServerError{
		Message:    message,
		Code:       cfg.code,
		Timestamp:  time.Now(),
		StackTrace: stackTrace,
	}
}

type config struct {
	code         string
	includeStack bool
}

type Option func(*config)

func WithCode(code string) Option {
	return func(c *config) {
		c.code = code
	}
}

func WithoutStack() Option {
	return func(c *config) {
		c.includeStack = false
	}
}

func captureStackTrace(skip int) string {
	var pcs [32]uintptr
	n := runtime.Callers(skip, pcs[:])
	frames := runtime.CallersFrames(pcs[:n])

	var sb strings.Builder
	for {
		frame, more := frames.Next()
		sb.WriteString(fmt.Sprintf("%s\n\t%s:%d\n", frame.Function, frame.File, frame.Line))
		if !more {
			break
		}
	}
	return sb.String()
}
