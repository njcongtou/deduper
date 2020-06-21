package main

import (
	"context"
	"deduper"
	"fmt"
	"log"
	"net/http"
	"os"

	"k8s.io/client-go/rest"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
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

	var (
		clientset *kubernetes.Clientset
		err       error
		addrs     []string
	)

	// 1. creates the out of cluster config
	//if clientset, err = common.InitClient(); err != nil {
	//	fmt.Println(err)
	//	return
	//}

	// 2. creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}

	// 3. creates the clientset
	clientset, err = kubernetes.NewForConfig(config)

	configMap, err := clientset.CoreV1().ConfigMaps("default").Get(context.TODO(), "special-config", metav1.GetOptions{})
	if err != nil {
		panic(err.Error())
	}

	for _, value := range configMap.Data {
		fmt.Println("a value is : " + value)
		addrs = append(addrs, value)
	}

	fmt.Println("new constructed addresses are:")
	fmt.Println(addrs)

	podIP := os.Getenv("MY_POD_IP")
	cacheServer := podIP + ":8001"
	apierver := podIP + ":9999"

	group := createGroup()
	go startAPIServer(apierver, group)
	startCacheServer(cacheServer, []string(addrs), group)
}
