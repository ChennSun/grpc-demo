package main

import (
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc/balancer/roundrobin"
	_ "google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
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
	serverConfig := fmt.Sprintf(`{
    "LoadBalancingPolicy": "%s",
    "MethodConfig": [
        {
            "Name": [{"Service": "%s"}], 
			"Timeout": "5s",
            "RetryPolicy": {
                "MaxAttempts":5, 
                "InitialBackoff": "0.1s", 
                "MaxBackoff": "1s", 
                "BackoffMultiplier": 2.0,
                "RetryableStatusCodes": ["UNAVAILABLE", "CANCELLED", "DEADLINE_EXCEEDED"]
                }
        }
    ], 
    "HealthCheckConfig": {
        "ServiceName": "%s"
    }
}`, roundrobin.Name, hello.Hi_ServiceDesc.ServiceName, healthCheckService)
	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr, grpc.WithInsecure(), grpc.WithDefaultServiceConfig(serverConfig))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := hello.NewHiClient(conn)
	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()
	r, err := client.SayHello(ctx, &hello.HiUser{Name: *name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.GetUser())
	// 健康检查
	go func() {
		cli := grpc_health_v1.NewHealthClient(conn)
		// 传入的Service必须和SetServingStatus方法的服务名一致，服务名查询不到会抛error
		res, _ := cli.Watch(context.TODO(), &grpc_health_v1.HealthCheckRequest{
			Service: healthCheckService,
		})
		for {
			r, _ := res.Recv()
			if r == nil {
				continue
			}
			fmt.Println(r)
		}
	}()
}
