package appgroup_test

import (
	"context"
	"testing"
	"time"

	"github.com/brunograsselli/appgroup"
)

func TestWithContext(t *testing.T) {
	ctx := context.Background()

	g, ctx := appgroup.WithContext(ctx)

	for i := 0; i < 3; i++ {
		g.Go(func() {})
	}

	g.Wait()

	want, got := context.Canceled, ctx.Err()
	if got != want {
		t.Errorf("unexpected error, got %v, want %v", got, want)
	}
}

func TestWithContextTimeout(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	g, _ := appgroup.WithContext(ctx)

	closeCh := make(chan struct{})

	g.Go(func() {
		// Wait forever
		<-closeCh
	})

	doneCh := make(chan struct{})

	go func() {
		defer close(doneCh)
		g.Wait(appgroup.WithShutdownTimeout(100 * time.Millisecond))
	}()

	cancel()
	<-doneCh
}
