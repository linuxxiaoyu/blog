package main

import (
	"log"
	"net"

	pb "github.com/linuxxiaoyu/blog/api"
	"github.com/linuxxiaoyu/blog/internal/setting"
	"google.golang.org/grpc"
)

func main() {
	setting.InitDB()
	setting.InitCache()
	lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterAccountServer(s, &service{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
