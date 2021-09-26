package main

import (
	"context"
	"fmt"
	"github.com/dairongpeng/gbase"
)

func main() {
	ctx := context.Background()
	LogDebug(ctx)
	// ConfigDebug(ctx)
}

func LogDebug(ctx context.Context) {
	gbase.Info(ctx).Msg("hello world")
	childCtx := gbase.AddLogValues(ctx, "name", "tom")
	gbase.Info(childCtx).Msg("print log")
}

func ConfigDebug(ctx context.Context) {
	name := gbase.Cfg().GetString("name")
	fmt.Println(name)
	port := gbase.Cfg().GetString("http.port")
	fmt.Println(port)
}
