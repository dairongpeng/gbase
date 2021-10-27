package main

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
)

func main() {
	group, _ := errgroup.WithContext(context.Background())
	for i := 0; i < 5; i++ {
		group.Go(func() error {
			return nil
		})
	}
	if err := group.Wait(); err != nil {
		fmt.Println(err)
	}
}
