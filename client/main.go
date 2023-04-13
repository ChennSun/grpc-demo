package main

import (
	"context"
	"flag"
	"fmt"
	"hello"
	"log"
	"time"

	"google.golang.org/grpc"
)

const (
	defaultName        = "world"
	healthCheckService = "grpc.health.v1.Health"
)

var (
	addr = flag.String("addr", "localhost:9878", "the address to connect to")
	name = flag.String("name", defaultName, "Name to greet")
)

func main() {
	flag.Parse()
	// health check config
	serverConfig := fmt.Sprintf(`{"HealthCheckConfig": {"ServiceName": "%s"}}`, healthCheckService)
	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr, grpc.WithInsecure(), grpc.WithDefaultServiceConfig(serverConfig))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := hello.NewHiClient(conn)
	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := client.SayHello(ctx, &hello.HiUser{Name: *name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.GetUser())
}
