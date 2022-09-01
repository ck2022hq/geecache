package main

import (
	"flag"
	"fmt"
	"geecache"
	"log"
	"net/http"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
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
	peers := geecache.NewHTTPPool(addr)
	peers.Set(addrs...)
	log.Println("geecache is running at", addr)
	mux := http.NewServeMux()
	mux.Handle("/geecache", http.HandlerFunc(peers.ServeHTTP))
	log.Fatal(http.ListenAndServe(addr[7:], mux))
}
