// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"flag"

	"github.com/dairongpeng/gbase/zlog"
)

var (
	h bool

	level  int
	format string
)

func main() {
	flag.BoolVar(&h, "h", false, "Print this help.")
	flag.IntVar(&level, "l", 0, "Log level.")
	flag.StringVar(&format, "f", "console", "log output format.")

	flag.Parse()

	if h {
		flag.Usage()

		return
	}

	// logger配置
	opts := &zlog.Options{
		Level:            "debug",
		Format:           "console",
		EnableColor:      true,
		DisableCaller:    true,
		OutputPaths:      []string{"test.log", "stdout"},
		ErrorOutputPaths: []string{"error.log"},
	}
	// 初始化全局logger
	zlog.Init(opts)
	defer zlog.Flush()

	// WithValues使用
	lv := zlog.WithValues("X-Request-ID", "7a7b9f24-4cae-4b2a-9464-69088b45b904")

	// Context使用
	lv.Infof("Start to call pirntString function")
	ctx := lv.WithContext(context.Background())
	pirntString(ctx, "World")
}

func pirntString(ctx context.Context, str string) {
	lc := zlog.FromContext(ctx)
	lc.Infof("Hello %s", str)
}
