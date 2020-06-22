package main

import (
	"context"
	"deduper"
	"deduper/common"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"k8s.io/client-go/rest"

	v1 "k8s.io/api/core/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

var numOfEntries int

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

func startCacheServer(addr string, addrs []string, group *deduper.Group, clientset *kubernetes.Clientset, configMap *v1.ConfigMap) {
	peers := deduper.NewHTTPPool(addr)
	peers.Set(addrs...)
	group.RegisterPeers(peers)
	go checkConfigMapChange(peers, clientset)
	log.Println("dedupercache is running at", addr)
	log.Fatal(http.ListenAndServe(addr, peers))
}

func checkConfigMapChange(peers *deduper.HTTPPool, clientset *kubernetes.Clientset) {

	// check for configmap updates every 10s
	for {
		time.Sleep(10 * time.Second)

		configMap, err := clientset.CoreV1().ConfigMaps("default").Get(context.TODO(), "special-config", metav1.GetOptions{})
		if err != nil {
			panic(err.Error())
		}

		var addrs []string

		if len(configMap.Data) != numOfEntries {
			for _, v := range configMap.Data {
				addrs = append(addrs, v)
			}
			peers.Set(addrs...)
			log.Println("new set of  peers are set successfully !")
			log.Println(configMap.Data)
			numOfEntries = len(configMap.Data)
		} else {
			//log.Println("configMap is not updated, no need to reset peers!")
		}

	}

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
	if clientset, err = common.InitClient(); err != nil {
		// 2. change to in cluster config if this is not local dev
		config, err := rest.InClusterConfig()
		if err != nil {
			panic(err.Error())
		}
		clientset, err = kubernetes.NewForConfig(config)
	}

	configMap, err := clientset.CoreV1().ConfigMaps("default").Get(context.TODO(), "special-config", metav1.GetOptions{})
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("new constructed addresses are:")
	fmt.Println(addrs)

	// add myself to configmap if it does not exist

	podIP := os.Getenv("MY_POD_IP")
	cacheServer := podIP + ":8001"
	apierver := podIP + ":9999"

	writeSelftoConfigMap(configMap, cacheServer, clientset, podIP)

	group := createGroup()
	go startAPIServer(apierver, group)
	startCacheServer(cacheServer, []string(addrs), group, clientset, configMap)
}

func writeSelftoConfigMap(configMap *v1.ConfigMap, cacheServer string, clientset *kubernetes.Clientset, podIP string) {
	if podIP == "" {
		podIP = "localhost"
	}

	exist := false
	for key, _ := range configMap.Data {
		if podIP == key {
			exist = true
		}
	}
	// Update configmap
	if !exist {
		fmt.Println("Updating configmap with new entry!")
		configMap.Data[podIP] = cacheServer
		fmt.Println(podIP + ", " + cacheServer)

		_, err := clientset.CoreV1().ConfigMaps("default").Update(context.TODO(), configMap, metav1.UpdateOptions{})
		if err != nil {
			fmt.Println("error : " + err.Error())
		}
		numOfEntries = 1
	}
}
