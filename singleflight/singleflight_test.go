package singleflight

import (
	"log"
	"sync"
	"testing"
	"time"
)

func TestQuery(t *testing.T) {
	g := Group{}
	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		go func() {
			wg.Add(1)
			defer wg.Done()
			v, err := g.Query("key", func() (interface{}, error) {
				time.Sleep(5 * time.Second)
				return "bar", nil
			})
			log.Println("res:", v, err)
		}()
	}
	time.Sleep(4 * time.Second)
	wg.Wait()
}

func TestDo(t *testing.T) {
	g := Group{}
	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		go func() {
			wg.Add(1)
			defer wg.Done()
			v, err := g.Do("key", func() (interface{}, error) {
				time.Sleep(5 * time.Second)
				return "bar", nil
			})
			log.Println("res:", v, err)
		}()
	}
	time.Sleep(4 * time.Second)
	wg.Wait()
}
