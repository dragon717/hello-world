package main

import (
	"Test/protocol/cs"
	"context"
	"log"
	"net"
	"sync"

	"google.golang.org/grpc"
)

// 在 server/core目录下 升级指定版本的grpc包
// go get google.golang.org/grpc
// 查看grpc包版本
// go list -m google.golang.org/grpc

type userService struct {
	cs.UnimplementedUserServiceServer
}

// 全局保存所有活跃的NotifyStream连接
var notifyStreams sync.Map // key: user_id, value: cs.UserService_NotifyStreamServer

func (s *userService) GetUser(ctx context.Context, req *cs.C2S_GetUser) (*cs.S2C_GetUser, error) {
	log.Printf("Received request for user ID: %d", req.Id)
	user := &cs.User{
		Id: req.Id,
		// 其他字段填充逻辑
	}
	return &cs.S2C_GetUser{User: user}, nil
}

func (s *userService) NotifyStream(stream cs.UserService_NotifyStreamServer) error {
	var userID uint64
	for {
		req, err := stream.Recv()
		if err != nil {
			log.Printf("NotifyStream Recv error: %v", err)
			if userID != 0 {
				notifyStreams.Delete(userID)
			}
			return err
		}
		userID = req.GetUserId()
		notifyStreams.Store(userID, stream)
		log.Printf("NotifyStream connected: user_id=%d", userID)
		// 可以处理客户端发来的消息（如心跳等）
	}
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
