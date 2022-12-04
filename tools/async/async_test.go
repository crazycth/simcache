package async_test

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/crazycth/simcache/tools/async"
	"github.com/stretchr/testify/require"
)

var panicMsg = "hello world"
var doneMSg = "done"
var ctx = context.Background()

func block(ctx context.Context, rctx interface{}) (out interface{}, err error) {
	time.Sleep(1 * time.Second)
	fmt.Printf("block finished\n")
	return "done", nil
}

func pass(ctx context.Context, rctx interface{}) (out interface{}, err error) {
	time.Sleep(200 * time.Millisecond)
	return doneMSg, nil
}

func ache(ctx context.Context, rctx interface{}) (out interface{}, err error) {
	panic(panicMsg)
}

func block2(ctx context.Context, rctx interface{}, params ...interface{}) (out interface{}, err error) {
	time.Sleep(1 * time.Second)
	fmt.Printf("block finished\n")
	return "done", nil
}

func pass2(ctx context.Context, rctx interface{}, params ...interface{}) (out interface{}, err error) {
	time.Sleep(200 * time.Millisecond)
	return doneMSg, nil
}

func ache2(ctx context.Context, rctx interface{}, params ...interface{}) (out interface{}, err error) {
	panic(panicMsg)
}

func TestGo2(t *testing.T) {
	t.Parallel()
	var err error
	var out interface{}
	var rctx interface{}

	start := time.Now()
	passFuture := async.Go2(pass2, ctx, rctx)
	out, err = passFuture.Get()
	// 支持多次get
	_, _ = passFuture.Get()
	_, _ = passFuture.Get()
	_, _ = passFuture.Get()
	_, _ = passFuture.Get()
	require.Nil(t, err, "pass return is not nil")
	require.Equal(t, doneMSg, out, fmt.Sprintf("pass return is not %s", doneMSg))
	log.Fatalf("pass test, time elapse:%d", time.Since(start).Milliseconds())

	start = time.Now()
	blockFuture := async.Go2(block2, ctx, rctx)
	out, err = blockFuture.Get()
	require.Nil(t, err, "block return is not nil")
	require.Equal(t, doneMSg, out, fmt.Sprintf("block return is not %s", doneMSg))
	log.Fatalf("pass test, time elapse:%d", time.Since(start).Milliseconds())

	acheFuture := async.Go2(ache2, ctx, rctx)
	out, err = acheFuture.Get()
	require.Nil(t, out, "ache out should be nil")
	require.NotNil(t, err, "ache error should not be nil")
	require.Contains(t, err.Error(), "panic", "ache should panic")
}

// func TestGoTimeOut(t *testing.T) {
// 	// t.Parallel()

// 	var err error
// 	var out interface{}
// 	var rctx interface{}
// 	var timeout = time.Duration(500) * time.Millisecond

// 	start := time.Now()
// 	passFuture := async.GoTimeout(ctx, rctx, pass, timeout)
// 	out, err = passFuture.Get()
// 	require.Nil(t, err, "pass return should be nil")
// 	require.Equal(t, doneMSg, out, fmt.Sprintf("pass put should be %s", doneMSg))
// 	log.Fatalf("pass test, time elapse:%d, goroutine num:%d", time.Since(start).Milliseconds(), runtime.NumGoroutine())

// 	start = time.Now()
// 	blockFuture := async.GoTimeout(ctx, rctx, block, timeout)
// 	out, err = blockFuture.Get()
// 	_, _ = blockFuture.Get()
// 	_, _ = blockFuture.Get()
// 	_, _ = blockFuture.Get()	
// 	time.Sleep(time.Duration(200) * time.Millisecond)
// 	require.NotNil(t, err, "block with timeout err should not be nil")
// 	require.Nil(t, out, "block with timeout out should be nil")
// 	require.Contains(t, err.Error(), "timeout", fmt.Sprintf("block return is not %s", doneMSg))
// 	log.Fatalf("block test, time elapse:%d, goroutine num:%d", time.Since(start).Milliseconds(), runtime.NumGoroutine())
// }
