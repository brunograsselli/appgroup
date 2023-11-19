package appgroup

import (
	"context"
	"sync"
	"time"
)

var defaultShutdownTimeout = 30 * time.Second

type Group interface {
	Wait()
	Go(func())
}

type groupImpl struct {
	cancel  func()
	wg      sync.WaitGroup
	closeCh chan struct{}
	once    sync.Once
}

type waitConfig struct {
	shutdownTimeout time.Duration
}

type WaitOption func(*waitConfig)

func WithShutdownTimeout(t time.Duration) WaitOption {
	return func(o *waitConfig) {
		o.shutdownTimeout = t
	}
}

func WithContext(ctx context.Context) (*groupImpl, context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	g, ctx := &groupImpl{cancel: cancel, closeCh: make(chan struct{})}, ctx

	// Initiate shutdown if context is done
	g.Go(func() {
		<-ctx.Done()
	})

	return g, ctx
}

func (g *groupImpl) Wait(opts ...WaitOption) {
	// Load config
	cfg := &waitConfig{
		shutdownTimeout: defaultShutdownTimeout,
	}

	for _, opt := range opts {
		opt(cfg)
	}

	// Wait until shutdown is initiated
	<-g.closeCh

	// Wait for WaitGroup with timeout
	c := make(chan struct{})

	go func() {
		defer close(c)
		g.wg.Wait()
	}()

	select {
	case <-c:
	case <-time.After(cfg.shutdownTimeout):
	}

	// Cancel context
	g.cancel()
}

func (g *groupImpl) Go(f func()) {
	g.wg.Add(1)

	go func() {
		f()
		g.once.Do(g.initiateShutdown)
		g.wg.Done()
	}()
}

func (g *groupImpl) initiateShutdown() {
	close(g.closeCh)
	g.cancel()
}
