package main

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"hello"
	"log"
	"net"
)

const healthCheckService = "grpc.health.v1.Health"

type server struct {
	hello.UnimplementedHiServer
}

func (s *server) SayHello(ctx context.Context, user *hello.HiUser) (*hello.HiReply, error) {
	return &hello.HiReply{
		Message: "hello",
		User:    user,
	}, nil
}

func main() {
	// 开启一个tcp服务器
	l, err := net.Listen("tcp", "localhost:9878")
	if err != nil {
		log.Fatalln("tcp监听失败")
	}
	// 注册grpc服务
	s := grpc.NewServer()
	hello.RegisterHiServer(s, &server{})
	// 注册健康检查服务
	h := health.NewServer()
	// 可以修改对应服务的状态
	h.SetServingStatus(healthCheckService, grpc_health_v1.HealthCheckResponse_SERVING)
	grpc_health_v1.RegisterHealthServer(s, h)

	log.Printf("server listening at %v", l.Addr())
	if err := s.Serve(l); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
