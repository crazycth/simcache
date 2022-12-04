package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/crazycth/simcache/simcache"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func createGroup() *simcache.Group {
	return simcache.NewGroup("scores", 2<<10, simcache.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))
}

func startCacheServer(addr string, addrs []string, cache *simcache.Group) {
	//1. build httppool according to addr
	pool := simcache.NewHTTPPool(addr)

	//2. build several peers
	pool.Set(addrs...)

	//3. cache -> httppool
	cache.RegisterPeers(pool)

	log.Println("cache is running at", addr)
	log.Fatal(http.ListenAndServe(addr[7:], pool))
}

func startAPIServer(apiAddr string, addrs []string, gee *simcache.Group) {
	//1. init peers
	pool := simcache.NewHTTPPool(apiAddr)
	pool.Set(addrs...)
	gee.RegisterPeers(pool)

	//2. api server
	http.Handle("/api", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			name := r.URL.Query().Get("name")
			key := r.URL.Query().Get("key")
			log.Printf("[server] get query name : %s , key : %s", name, key)
			view, err := gee.QueryPeer(name, key)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/octet-stream")
			w.Write(view.ByteSlice())
		}))

	log.Println("fontend server is running at", apiAddr)
	log.Fatal(http.ListenAndServe(apiAddr[7:], nil))
}

func main() {
	var port int
	var api bool
	flag.IntVar(&port, "port", 8001, "simacache server port")
	flag.BoolVar(&api, "api", false, "start a api server?")
	flag.Parse()

	apiAddr := "http://localhost:9999"
	addrMap := map[int]string{
		8001: "http://localhost:8001",
		8002: "http://localhost:8002",
		8003: "http://localhost:8003",
	}

	var addrs []string
	for _, v := range addrMap {
		addrs = append(addrs, v)
	}

	cc := createGroup()
	if api {
		startAPIServer(apiAddr, addrs, cc)
	} else {
		startCacheServer(addrMap[port], addrs, cc)
	}
}
