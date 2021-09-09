package main

import (
	"fmt"
	"github.com/gopher-core/base"
)

func main() {
	cfg := base.Config()
	fmt.Println(cfg.Get("mysql.web.addr"))
}
