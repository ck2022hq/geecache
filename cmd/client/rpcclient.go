package main

import (
	"context"
	"flag"
	"log"
	"time"

	"google.golang.org/grpc"

	pb "geecache/geecachepb"
)

const (
	address = "localhost:11111"
)

var key = flag.String("key", "Jack", "key")

func main() {
	flag.Parse()
	//建立链接
	conn, err := grpc.Dial(address, grpc.WithInsecure())

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()

	cli := pb.NewGroupCacheClient(conn)

	// 设定请求超时时间 3s
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	req := &pb.Request{Group: "haha", Key: *key}
	rsp, err := cli.Get(ctx, req)
	if err != nil {
		log.Printf("rpc get fail, req=%s, err=%v", req, err)
	} else {
		log.Printf("rpc get success, rsp=%s", rsp)
	}
}
