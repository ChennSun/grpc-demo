package main

import (
	"context"
	"google.golang.org/grpc"
	"hello"
	"log"
	"net"
)

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
	log.Printf("server listening at %v", l.Addr())
	if err := s.Serve(l); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
