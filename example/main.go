package main

import (
	"context"
	"github.com/dairongpeng/gbase"
)

func main() {
	ctx := context.Background()
	gbase.Info(ctx).Str("hello", "world").Msg("test log")
}
