package geecache

import (
	"context"
	"fmt"
	"log"

	pb "geecache/geecachepb"
)

// implement GroupCache service
type GroupCacheService struct {
}

func (service *GroupCacheService) Get(ctx context.Context, req *pb.Request) (*pb.Response, error) {
	log.Printf("incoming rpc request, req=%s", req)

	rsp := &pb.Response{}
	if len(req.Group) == 0 || len(req.Key) == 0 {
		err := fmt.Errorf("invalid request, group or key empty")
		return rsp, err
	}

	val, err := Pool.Process(req.Group, req.Key)
	if err == nil {
		rsp.Value = val.ByteSlice()
	}

	return rsp, err
}
