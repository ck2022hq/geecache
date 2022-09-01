go build -o server
./server -port=8001 &
./server -port=8002 &
./server -port=8003 &
./server -port=9999 &