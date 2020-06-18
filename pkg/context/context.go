package context

import (
	"context"
	"sync"
)

type WaitStopContext struct {
	ctx        context.Context
	cancelFunc context.CancelFunc
	wg         sync.WaitGroup
}

func NewWaitStopContext() *WaitStopContext {
	ctx, cancel := context.WithCancel(context.Background())
	return &WaitStopContext{
		ctx:        ctx,
		cancelFunc: cancel,
	}
}

func (c *WaitStopContext) Cancel() {
	c.cancelFunc()
	c.wg.Wait()
}

func (c *WaitStopContext) Add(delta int) {
	c.wg.Add(delta)
}

func (c *WaitStopContext) Finish() {
	c.wg.Done()
}

func (c *WaitStopContext) Done() <-chan struct{} {
	return c.ctx.Done()
}
