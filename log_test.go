package gbase

import (
	"context"
	"testing"
)

func TestLogCtxKey(t *testing.T) {
	ctx := context.Background() // is emptyCtx and root ctx
	cm := map[string]string{"name": "tom", "age": "18"}
	newCtx := context.WithValue(ctx, "logCtxKey", cm)
	// value := newCtx.Value("logCtxKey")
	// m := value.(map[string]string)
	// gbase.Error(newCtx).Str("hello", "world").Msg("test log")
	Info(newCtx).Int("length", 100).Str("hello", "world").Msg("test log ctx key")
}
