package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	_ "github.com/accelforce/lemokon/db"
	"github.com/accelforce/lemokon/kon"
	"github.com/accelforce/lemokon/lemo"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	wg := new(sync.WaitGroup)

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	go func() {
		<-sc
		cancel()
	}()

	wg.Add(2)
	go func() {
		if err := lemo.Lemo(ctx); err != nil {
			fmt.Printf("Error happened in Lemo: %s\n", err)
			cancel()
		}
		wg.Done()
	}()
	go func() {
		if err := kon.Kon(ctx); err != nil {
			fmt.Printf("Error happened in Kon: %s\n", err)
			cancel()
		}
		wg.Done()
	}()
	wg.Wait()
}
