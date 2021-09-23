package gbase

import (
	"context"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"os"
	"runtime"
	"strings"
	"time"
)

var Logger zerolog.Logger
var nilCtx = context.Background()
var ctxLogKey ctxLogKeyType = "logFields"

type ctxLogKeyType string
type logLevel string

const (
	TRACE logLevel = "TRACE"
	DEBUG logLevel = "DEBUG"
	WARN  logLevel = "WARN"
	ERROR logLevel = "ERROR"
)

func init() {
	zerolog.TimeFieldFormat = time.RFC3339Nano
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack // Error().Stack().Err(err).Msg("") will print err stack
	zerolog.TimestampFieldName = "t"
	zerolog.LevelFieldName = "level"
	zerolog.MessageFieldName = "msg"
	Logger = zerolog.New(os.Stdout).With().Logger().Level(zerolog.InfoLevel).Hook(LogHook{})
	level := Viper().GetString("log.level")
	switch logLevel(level) {
	case TRACE:
		Logger = Logger.Level(zerolog.TraceLevel)
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	case DEBUG:
		Logger = Logger.Level(zerolog.DebugLevel)
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case WARN:
		Logger = Logger.Level(zerolog.WarnLevel)
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case ERROR:
		Logger = Logger.Level(zerolog.ErrorLevel)
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	default:
		Logger = Logger.Level(zerolog.DebugLevel)
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	log.Logger = Logger
}

type LogHook struct{}

func (LogHook) Run(event *zerolog.Event, level zerolog.Level, message string) {
	// TODO add metric for log detail info
}

// WithLogContext let ctx k/v package to str log env
func WithLogContext(ctx context.Context, e *zerolog.Event) *zerolog.Event {
	if ctx == nilCtx {
		return e
	}
	logFields := fromCtxLogItems(ctx)
	if len(logFields) == 0 {
		return e
	}

	for k, v := range logFields {
		e = e.Str(k, v)
	}
	return e
}

// fromCtxLogItems will decode ctx than get k v
func fromCtxLogItems(ctx context.Context) map[string]string {
	raw := ctx.Value(ctxLogKey)
	if raw == nil {
		return map[string]string{}
	}
	return raw.(map[string]string)
}

// appendEvents decide log user call or not
func appendEvents(event *zerolog.Event, addCaller bool) *zerolog.Event {
	event.Timestamp()
	if addCaller {
		_, file, line := funcFileLine("github.com/dairongpeng/gbase")
		event.Str("caller", fmt.Sprintf("%s:%d", file, line))
	}
	return event
}

// funcFileLine will build caller item. The aim is to send to upper layer App
func funcFileLine(excludePKG string) (string, string, int) {
	const depth = 8
	var pcs [depth]uintptr
	n := runtime.Callers(3, pcs[:])
	ff := runtime.CallersFrames(pcs[:n])

	var fn, file string
	var line int
	for {
		f, ok := ff.Next()
		if !ok {
			break
		}
		fn, file, line = f.Function, f.File, f.Line
		if !strings.Contains(fn, excludePKG) {
			break
		}
	}

	if ind := strings.LastIndexByte(fn, '/'); ind != -1 {
		fn = fn[ind+1:]
	}

	return fn, file, line
}

// Debug will get log.event with DEBUG level
func Debug(ctx context.Context) *zerolog.Event {
	return WithLogContext(ctx, appendEvents(log.Debug(), false))
}

// DebugWithoutCtx will get log.event with DEBUG level without context
func DebugWithoutCtx() *zerolog.Event {
	return Debug(nilCtx)
}

// Info will get log.event with INFO level
func Info(ctx context.Context) *zerolog.Event {
	return WithLogContext(ctx, appendEvents(log.Info(), false))
}

// InfoWithoutCtx will get log.event with INFO level without context
func InfoWithoutCtx() *zerolog.Event {
	return Info(nilCtx)
}

// Warn will get log.event with WARN level
func Warn(ctx context.Context) *zerolog.Event {
	return WithLogContext(ctx, appendEvents(log.Warn(), true))
}

// WarnWithoutCaller will get log.event with WARN level no caller
func WarnWithoutCaller(ctx context.Context) *zerolog.Event {
	return WithLogContext(ctx, appendEvents(log.Warn(), false))
}

// WarnWithoutCtx will get log.event with WARN level without ctx
func WarnWithoutCtx() *zerolog.Event {
	return Warn(nilCtx)
}

// Error will get log.event with ERROR level without ctx
func Error(ctx context.Context) *zerolog.Event {
	return WithLogContext(ctx, appendEvents(log.Error(), true))
}

// Err will get log.event with ERROR level
func Err(ctx context.Context, err error) *zerolog.Event {
	return WithLogContext(ctx, appendEvents(log.Err(err), true))
}

// ErrorWithoutCaller will get log.event with ERROR level no caller
func ErrorWithoutCaller(ctx context.Context) *zerolog.Event {
	return WithLogContext(ctx, appendEvents(log.Error(), false))
}

// ErrorWithoutCtx will get log.event with ERROR level without ctx
func ErrorWithoutCtx() *zerolog.Event {
	return Error(nilCtx)
}

// Fatal will get log.event with FATAL level
func Fatal(ctx context.Context) *zerolog.Event {
	return WithLogContext(ctx, appendEvents(log.Fatal(), true))
}

// FatalWithoutCtx will get log.event with FATAL level without ctx
func FatalWithoutCtx() *zerolog.Event {
	return Fatal(nilCtx)
}

// WarnErr will get log.event with WARN and contains a ERROR level log
func WarnErr(ctx context.Context, err error) {
	Warn(ctx).Err(err).Send()
}
