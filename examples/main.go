package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"
	"time"

	"github.com/brunograsselli/appgroup"
)

func main() {
	ctx := context.Background()

	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	g, ctx := appgroup.WithContext(ctx)

	g.Go(func() {
		ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
		defer cancel()

		<-ctx.Done() // simulating a long running listener that takes context into account

		fmt.Println("ending first go routine")
	})

	g.Go(func() {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		<-ctx.Done() // simulating a long running listener that takes context into account

		fmt.Println("ending second go routine")
	})

	g.Wait(appgroup.WithShutdownTimeout(10 * time.Second))

	fmt.Println("done")
}
