#!/bin/bash
trap "rm server;kill 0" EXIT

go build -o server
./server -port=8001 -api=0 &
./server -port=8002 -api=0 &
./server -port=8003 -api=0 &
./server -port=8004 -api=1 &

sleep 2
echo ">>> start test"
curl "http://localhost:9999/api?key=Tom&name=scores" &
curl "http://localhost:8002/_richardcache/scores/Tom" &
# curl "http://localhost:9999/api?key=Tom&name=scores" &


wait
