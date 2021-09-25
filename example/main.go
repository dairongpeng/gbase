package main

import (
	"context"
	"github.com/dairongpeng/gbase"
)

func main() {
	ctx := context.Background()
	gbase.Info(ctx).Msg("hello world")
	childCtx := gbase.AddLogValues(ctx, "name", "tom")
	gbase.Info(childCtx).Msg("print log")
}
