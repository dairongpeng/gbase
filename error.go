package gbase

import (
	"fmt"
	"os"
)

// TODO 使用 errors及错误状态码进行再次封装

// CheckErr prints the msg with the prefix 'Error:' and exits with error code 1. If the msg is nil, it does nothing.
func CheckErr(msg interface{}) {
	if msg != nil {
		_, _ = fmt.Fprintln(os.Stderr, "Error:", msg)
		os.Exit(1)
	}
}
