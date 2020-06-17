# deduper

## New Features

1. local cache with Get() function
2. lru 
3. http server
4. Group search
5. docker support
6. localhost port 8080
7. URL : http://localhost:8080/_dedupercache/scores/Tom
8. consistent hasher

    pattern: default_base_url/group/key

## Dev commands:


  docker stop `docker ps -a -q`

  docker build -t deduper .

  docker run -p 8080:8080 -d deduper

  curl http://localhost:8080/_dedupercache/scores/Tom
  
 ## Dev Goal
 
 ### Current Goal
 
  1. add distributed nodes support (local ports first)
  2. create 3 statefulset pods, find IPs. find a way to register them on hashring
 
 ### Final Goal
 
   k8s statefulset pods are allocated on a hashring. Pod is able to send msg to others based on hashring.
   k8s operator: ?
   1. When initilizing statefulset pods, all pods' IPs should be registered on the hashring. 
   2. When adding new pod(s) (scale up cluster), newly added pod's IP should be registered on all pods.
   3. When removing existing pod(s) (scale up cluster), the pod's IP should be removed on all pods.
   4. Properly handing statefulset pod restarting, not scale up or down. Just restarted. ?
   5. Properly handing extreme case: kvm (node) failure, pods are scheduled on other nodes.
