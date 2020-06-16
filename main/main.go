package main

import (
	"deduper2"
	"fmt"
	"log"
	"net/http"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func main() {
	deduper2.NewGroup("scores", 2<<10, deduper2.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))

	addr := ":8080"
	peers := deduper2.NewHTTPPool(addr)
	log.Println("dedupercache is running at", addr)
	log.Fatal(http.ListenAndServe(addr, peers))

}
