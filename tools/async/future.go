package async

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type FutureOut struct {
	out interface{}
	err error
}

type Future struct {
	mutex     sync.Mutex
	c         chan *FutureOut
	cOut      *FutureOut
	timeout   time.Duration
	cancelCtx context.Context
	cancel    context.CancelFunc
}

func (f *Future) getDirect() (interface{}, error) {
	if f.cOut == nil {
		f.mutex.Lock()
		defer f.mutex.Unlock()

		if f.cOut != nil {
			return f.cOut.out, f.cOut.err
		}
		f.cOut = <-f.c
	}

	return f.cOut.out, f.cOut.err
}

func (f *Future) Get() (interface{}, error) {
	// cancel the child goroutines
	if f.cancel == nil || f.cancelCtx == nil || f.timeout == 0 {
		return f.getDirect()
	}

	if f.cOut == nil {
		f.mutex.Lock()
		defer f.mutex.Unlock()
		if f.cOut != nil {
			return f.cOut.out, f.cOut.err
		}

		defer f.cancel()
		select {
		case f.cOut = <-f.c:
		// cancelCtx will be readable when current goroutine exit or cancel executed
		case <-f.cancelCtx.Done():
			err := fmt.Errorf("[core][tools][async][GetTimeout] future execute timeout:%d, err:%s", f.timeout.Milliseconds(), f.cancelCtx.Err().Error())
			f.cOut = &FutureOut{out: nil, err: err}
		}
	}

	return f.cOut.out, f.cOut.err
}
