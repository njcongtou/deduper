package main

import (
	"deduper"
	"fmt"
	"log"
	"net/http"
	"os"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func createGroup() *deduper.Group {
	return deduper.NewGroup("scores", 2<<10, deduper.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))
}

func startCacheServer(addr string, addrs []string, group *deduper.Group) {
	peers := deduper.NewHTTPPool(addr)
	peers.Set(addrs...)
	group.RegisterPeers(peers)
	log.Println("dedupercache is running at", addr)
	log.Fatal(http.ListenAndServe(addr, peers))
}

func startAPIServer(apiAddr string, group *deduper.Group) {
	http.Handle("/api", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			key := r.URL.Query().Get("key")
			view, err := group.Get(key)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/octet-stream")
			w.Write(view.ByteSlice())

		}))
	log.Println("fontend server is running at", apiAddr)
	log.Fatal(http.ListenAndServe(apiAddr, nil))

}

func main() {

	addrMap := map[string]string{
		"pod1": "172.17.0.2:8001",
		"pod2": "172.17.0.3:8001",
		"pod3": "172.17.0.10:8001",
	}

	var addrs []string
	for _, v := range addrMap {
		addrs = append(addrs, v)
	}

	podIP := os.Getenv("MY_POD_IP")
	cacheServer := podIP + ":8001"
	apierver := podIP + ":9999"

	group := createGroup()
	go startAPIServer(apierver, group)
	startCacheServer(cacheServer, []string(addrs), group)
}
