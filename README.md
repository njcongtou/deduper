# deduper

## Progress

1. local cache with Get() function
2. lru 
3. http server
4. Group search
5. docker support
6. localhost port 8080
7. URL : http://localhost:8080/_dedupercache/scores/Tom
8. consistent hasher
9. add distributed nodes support (tested using local ports 8001, 8002, 8003), api server on port 9999

    pattern: default_base_url/group/key
    
10. using Deployment instead of Statefulsets, since it does not support injecting pod ip env variable.
11. support multiple k8s pods exchanging keys.
12. support configmap:
    hardcoded 3 pods in the configmap
    read them from main.go

## Dev commands:

  docker build -t jimwallet/deduper:v2 . && docker push jimwallet/deduper:v2

  curl http://172.17.0.2:8001/_dedupercache/scores/Tom
  
  curl http://172.17.0.2:9999/api?key=Tom
  
  ### Local Testing In cluster, get from configmap, which requires permission.
  kubectl create clusterrolebinding default-admin --clusterrole cluster-admin --serviceaccount=default:default
  https://stackoverflow.com/questions/52954810/kubespray-dashboard-warning-forbidden-popups
  
 ## Dev Goal
 
 ### TODO
 
    * store membership in configmap, and a goroutine peridically checking it.
 
 ### Final Goal
 
   k8s statefulset pods are allocated on a hashring. Pod is able to send msg to others based on hashring.
   k8s operator: ?
   1. When initilizing statefulset pods, all pods' IPs should be registered on the hashring. 
   2. When adding new pod(s) (scale up cluster), newly added pod's IP should be registered on all pods.
   3. When removing existing pod(s) (scale up cluster), the pod's IP should be removed on all pods.
   4. Properly handing statefulset pod restarting, not scale up or down. Just restarted. ?
   5. Properly handing extreme case: kvm (node) failure, pods are scheduled on other nodes.
