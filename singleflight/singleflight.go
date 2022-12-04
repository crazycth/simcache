package singleflight

import (
	"context"
	"log"
	"sync"

	"github.com/crazycth/simcache/tools/async"
)

type call struct {
	wg  sync.WaitGroup
	val interface{}
	err error
}

type Group struct {
	M sync.Map
}

//FIXME 思考下能不能用异步tool来解决问题
func (g *Group) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
	if cI, ok := g.M.Load(key); ok {
		c := cI.(*call)
		c.wg.Wait() //如果该key已经有请求了，则等待
		return c.val, c.err
	}

	c := new(call)
	c.wg.Add(1) //发请求前加锁
	g.M.Store(key, c)

	c.val, c.err = fn()
	c.wg.Done()

	g.M.Delete(key)

	return c.val, c.err
}

//异步tool策略
func (g *Group) Query(key string, fn func() (interface{}, error)) (interface{}, error) {
	if futureI, ok := g.M.Load(key); ok {
		log.Fatalln("[singleflight][Query] hit future cache")
		future := futureI.(*async.Future)
		return future.Get()
	}

	log.Fatalln("[singleflight] start query")
	future := async.Go(fn, context.Background())
	g.M.Store(key, future)

	result, err := future.Get()
	g.M.Delete(key)
	return result, err
}
