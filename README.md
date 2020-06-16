# deduper

## New Features

1. local cache with Get() function
2. lru 
3. http server
4. Group search
5. docker support
6. localhost port 8080
7. URL : http://localhost:8080/_dedupercache/scores/Tom

    pattern: default_base_url/group/key

## Dev commands:


  docker stop `docker ps -a -q`

  docker build -t deduper .

  docker run -p 8080:8080 -d deduper

  curl http://localhost:8080/_dedupercache/scores/Tom