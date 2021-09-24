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

var (
	// Logger is global Logger
	Logger zerolog.Logger
	// initCtx Provides a init ctx
	initCtx = context.Background()
	// logCtxKey is ctx contains map key
	logCtxKey logCtxKeyType = "logCtxKey"
)

type (
	logCtxKeyType string
	logLevel      string
)

const (
	TRACE logLevel = "TRACE"
	DEBUG logLevel = "DEBUG"
	WARN  logLevel = "WARN"
	ERROR logLevel = "ERROR"
)

func init() {
	zerolog.TimeFieldFormat = time.RFC3339Nano
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack // Error().Stack().Err(err).Msg("") will print err stack
	zerolog.TimestampFieldName = "timestamp"
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

// WithLogContext returns Event. The event is already appends ctx kv
func WithLogContext(ctx context.Context, event *zerolog.Event) *zerolog.Event {
	// BackgroundCtx and TODOCtx are emptyCtx(root ctx). Compare BackgroundCtx and TODOCtx will return true
	if ctx == initCtx {
		return event
	}

	logCtxFields := fromCtxLogItems(ctx)
	if len(logCtxFields) == 0 {
		return event
	}

	for k, v := range logCtxFields {
		event = event.Str(k, v)
	}
	return event
}

// fromCtxLogItems returns the map from ctx contains kv
func fromCtxLogItems(ctx context.Context) map[string]string {
	raw := ctx.Value(string(logCtxKey))
	if raw == nil {
		return map[string]string{}
	}
	return raw.(map[string]string)
}

// appendEvents decides append caller or not
func appendEvents(event *zerolog.Event, addCaller bool) *zerolog.Event {
	event.Timestamp()
	if addCaller {
		_, file, line := funcFileLine("github.com/dairongpeng/gbase")
		event.Str("caller", fmt.Sprintf("%s:%d", file, line))
	}
	return event
}

// funcFileLine find cur err pkg be include
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

// Debug returns Event by debug level
func Debug(ctx context.Context) *zerolog.Event {
	return WithLogContext(ctx, appendEvents(log.Debug(), false))
}

// DebugWithoutCtx returns Event by debug level without ctx
func DebugWithoutCtx() *zerolog.Event {
	return Debug(initCtx)
}

// Info returns Event by info level
func Info(ctx context.Context) *zerolog.Event {
	return WithLogContext(ctx, appendEvents(log.Info(), false))
}

// InfoWithoutCtx returns Event by info level without ctx
func InfoWithoutCtx() *zerolog.Event {
	return Info(initCtx)
}

// Warn returns Event by warn level
func Warn(ctx context.Context) *zerolog.Event {
	return WithLogContext(ctx, appendEvents(log.Warn(), true))
}

// WarnWithoutCaller returns Event by warn level without caller
func WarnWithoutCaller(ctx context.Context) *zerolog.Event {
	return WithLogContext(ctx, appendEvents(log.Warn(), false))
}

// WarnWithoutCtx returns Event by warn level without ctx
func WarnWithoutCtx() *zerolog.Event {
	return Warn(initCtx)
}

// Error returns Event by error level
func Error(ctx context.Context) *zerolog.Event {
	return WithLogContext(ctx, appendEvents(log.Error(), true))
}

// Err returns Event by error level and print parameter err
func Err(ctx context.Context, err error) *zerolog.Event {
	return WithLogContext(ctx, appendEvents(log.Err(err), true))
}

// ErrorWithoutCaller returns Event by error level without caller
func ErrorWithoutCaller(ctx context.Context) *zerolog.Event {
	return WithLogContext(ctx, appendEvents(log.Error(), false))
}

// ErrorWithoutCtx returns Event by error level without ctx
func ErrorWithoutCtx() *zerolog.Event {
	return Error(initCtx)
}

// Fatal returns Event by fatal level
func Fatal(ctx context.Context) *zerolog.Event {
	return WithLogContext(ctx, appendEvents(log.Fatal(), true))
}

// FatalWithoutCtx returns Event by fatal level without ctx
func FatalWithoutCtx() *zerolog.Event {
	return Fatal(initCtx)
}

// WarnErr returns Event by warn level and print parameter err
func WarnErr(ctx context.Context, err error) {
	Warn(ctx).Err(err).Send()
}
