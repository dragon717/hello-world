package main

import (
	cs "Test/protocol"
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
)

// 在 server/core目录下 升级指定版本的grpc包
// go get google.golang.org/grpc
// 查看grpc包版本
// go list -m google.golang.org/grpc

type userService struct {
	cs.UnimplementedUserServiceServer
}

func (s *userService) GetUser(ctx context.Context, req *cs.C2S_GetUser) (*cs.S2C_GetUser, error) {
	log.Printf("Received request for user ID: %d", req.Id)
	user := &cs.User{
		Id: req.Id,
		// 其他字段填充逻辑
	}
	return &cs.S2C_GetUser{User: user}, nil
}

func startGRPCServer() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	cs.RegisterUserServiceServer(s, &userService{})
	log.Println("gRPC server running on port 50051...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
