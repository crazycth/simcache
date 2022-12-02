package simcache_test

import (
	"fmt"
	"log"
	"net/http"
)

type server int

func (h *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.Path)
	w.Write([]byte("Hello World!"))
}

func main() {
	var s server
	err := http.ListenAndServe("localhost:9999", &s)
	fmt.Println("err:", err)
}
