package appgroup_test

import (
	"context"
	"testing"

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
