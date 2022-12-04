package async

import (
	"context"
	"fmt"
	"runtime/debug"
	"time"
)

type FutureFunc2 func(ctx context.Context, rctx interface{}, params ...interface{}) (interface{}, error)
type FutureFn func() (interface{}, error)

func Go(fn FutureFn, ctx context.Context, params ...interface{}) *Future {
	c := make(chan *FutureOut, 1)
	go func() {
		var err error
		var ret interface{}

		defer func() {
			if r := recover(); r != nil {
				err := fmt.Errorf("[core][tools][async][Go] panic in async task err:%s, trace:%s", r, string(debug.Stack()))
				c <- &FutureOut{
					out: nil,
					err: err,
				}
			}
		}()

		ret, err = fn()
		c <- &FutureOut{
			out: ret,
			err: err,
		}
	}()
	return &Future{c: c}
}

// async with out timeout
func Go2(fn FutureFunc2, ctx context.Context, rctx interface{}, params ...interface{}) *Future {
	c := make(chan *FutureOut, 1)
	go func() {
		var err error
		var ret interface{}

		defer func() {
			if r := recover(); r != nil {
				err := fmt.Errorf("[core][tools][async][Go2] panic in async task err:%s, trace:%s", r, string(debug.Stack()))
				c <- &FutureOut{
					out: nil,
					err: err,
				}
			}
		}()

		ret, err = fn(ctx, rctx, params...)
		c <- &FutureOut{
			out: ret,
			err: err,
		}
	}()

	return &Future{c: c}
}

// async with timeout
// https://go.dev/blog/context http.Do implements a way that child goroutine will return directly
// https://bytedance.feishu.cn/wiki/wikcn4bVg5Fl2nXU7cAXUzjhZed#bJZ5Il
// goroutine will remain run until return by itself, except that when deadline reached when atempts to do the real thing
func GoTimeout2(fn FutureFunc2, ctx context.Context, timeout time.Duration, rctx interface{}, params ...interface{}) *Future {
	c := make(chan *FutureOut, 1)
	cancleCtx, cancel := context.WithTimeout(ctx, timeout)

	go func() {
		var err error
		var ret interface{}

		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("[core][tools][async][GoTimeout2] panic in async task err:%s, trace:%s", r, string(debug.Stack()))
				c <- &FutureOut{
					out: nil,
					err: err,
				}
			}
		}()

		deadline, _ := cancleCtx.Deadline()
		if !time.Now().Before(deadline) {
			c <- &FutureOut{
				out: nil,
				err: fmt.Errorf("[core][tools][async][GoTimeout2] deadline reached, timeout:%d, directly exit", timeout.Milliseconds()),
			}
			return
		}

		ret, err = fn(cancleCtx, rctx, params...)
		c <- &FutureOut{
			out: ret,
			err: err,
		}
	}()

	return &Future{
		c:         c,
		timeout:   timeout,
		cancelCtx: cancleCtx,
		cancel:    cancel,
	}
}
