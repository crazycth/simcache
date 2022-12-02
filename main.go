package main

import (
	"fmt"
	"github/richard003/simcache/simcache"
	"log"
	"net/http"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func main() {
	simcache.NewGroup("scores", 2<<10, simcache.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SolwDB] search key", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))
	addr := "localhost:9999"
	peers := simcache.NewHTTPPool(addr)
	log.Println("simcache is running at", addr)
	log.Fatal(http.ListenAndServe(addr, peers))
}
