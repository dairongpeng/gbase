package main

import (
	"fmt"
	_ "github.com/dairongpeng/gbase"
	"github.com/robfig/cron"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cron := cron.New()
	//nolint
	err := cron.AddFunc("0 * * * * ?", func() {
		fmt.Println("per s")
	})
	if err != nil {
		panic(err)
	}
	cron.Start()

	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}
