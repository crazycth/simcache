package singleflight

import (
	"context"
	"log"
	"sync"

	"github.com/crazycth/simcache/tools/async"
)

type Call struct {
	wg  sync.WaitGroup
	val interface{}
	err error
}

//结论: 不能用读写锁，否则对于同个key多个请求会同时读 --> 缓存击穿
type Group struct {
	mu      sync.Mutex
	FutureM map[string]*async.Future
	CallM   map[string]*Call
}

//TODO 能限住，但太多锁效率很低
func (g *Group) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
	g.mu.Lock()
	if g.CallM == nil {
		g.CallM = make(map[string]*Call)
	}
	if c, ok := g.CallM[key]; ok {
		log.Printf("[singleflight][Do] hit sync mutex")
		g.mu.Unlock()
		c.wg.Wait()
		return c.val, c.err
	}
	c := new(Call)
	c.wg.Add(1)
	g.CallM[key] = c
	g.mu.Unlock()

	c.val, c.err = fn()
	c.wg.Done()

	g.mu.Lock()
	delete(g.CallM, key)
	g.mu.Unlock()

	return c.val, c.err
}

//FIXME 思考有能提高效率的方法吗
func (g *Group) Query(key string, fn func() (interface{}, error)) (interface{}, error) {
	g.mu.Lock()
	if g.FutureM == nil {
		g.FutureM = make(map[string]*async.Future)
	}
	if future, ok := g.FutureM[key]; ok {
		g.mu.Unlock()
		log.Printf("[singleflight][Query] hit future cache key:%s", key)
		return future.Get()
	}

	log.Printf("[singleflight][Query] start query key:%s", key)

	future := async.Go(fn, context.Background())
	g.FutureM[key] = future
	g.mu.Unlock()

	result, err := future.Get()
	g.mu.Lock()
	delete(g.FutureM, key)
	g.mu.Unlock()
	return result, err
}
