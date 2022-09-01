package main

import (
	"flag"
	"fmt"
	"geecache"
	"log"
	"net"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "geecache/geecachepb"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

const (
	rpcPort = ":11111"
)

func startRpcService() {
	log.Println("start rpc service")
	lis, err := net.Listen("tcp", rpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterGroupCacheServer(grpcServer, &geecache.GroupCacheService{})
	reflection.Register(grpcServer)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func main() {
	geecache.DefaultGetter = geecache.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		})

	var port int
	flag.IntVar(&port, "port", 8001, "Geecache server port")
	flag.Parse()

	addrMap := map[int]string{
		8001: "http://localhost:8001",
		8002: "http://localhost:8002",
		8003: "http://localhost:8003",
	}

	var addrs []string
	for _, v := range addrMap {
		addrs = append(addrs, v)
	}

	addr, ok := addrMap[port]
	if !ok {
		addr = "http://localhost:9999"
	}
	geecache.Pool = geecache.NewHTTPPool(addr)
	geecache.Pool.Set(addrs...)
	log.Println("geecache is running at", addr)
	mux := http.NewServeMux()
	mux.Handle("/geecache", http.HandlerFunc(geecache.Pool.ServeHTTP))

	go startRpcService()

	log.Fatal(http.ListenAndServe(addr[7:], mux))
}
