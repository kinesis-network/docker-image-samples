package main

import (
	"context"
	"crypto/tls"
	"flag"
	"log"
	"net"

	"github.com/kinesis-network/go-greeter-client/greet"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	useTls = flag.Bool("tls", false, "Use TLS")
	target = flag.String("addr", "localhost:50051", "Server address")
	server = flag.Bool("s", false, "Server mode")
)

type greetServer struct {
	greet.UnimplementedGreeterServer
}

func (s *greetServer) SayHello(
	ctx context.Context,
	req *greet.HelloRequest,
) (*greet.HelloReply, error) {
	msg := "Hello " + req.Name
	log.Println(msg)
	return &greet.HelloReply{Message: msg}, nil
}

func main() {
	flag.Parse()
	var opts []grpc.DialOption
	if *useTls {
		creds := credentials.NewTLS(&tls.Config{
			InsecureSkipVerify: true,
		})
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(
			opts,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
	}

	if *server {
		var opts []grpc.ServerOption
		if *useTls {
			creds := credentials.NewTLS(&tls.Config{
				InsecureSkipVerify: true,
			})
			opts = append(opts, grpc.Creds(creds))
		}
		g := grpc.NewServer(opts...)
		greet.RegisterGreeterServer(g, &greetServer{})
		l, err := net.Listen("tcp", *target)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		log.Println("server listening on ", l.Addr().String())
		if err := g.Serve(l); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}

	g, err := grpc.NewClient(*target, opts...)
	if err != nil {
		log.Fatal("Error creating client: ", err)
	}
	defer g.Close()

	client := greet.NewGreeterClient(g)
	ctx := context.Background()
	resp, err := client.SayHello(ctx, &greet.HelloRequest{Name: "World"})
	if err != nil {
		log.Fatal("SayHello failed: ", err)
	}
	log.Println(resp.Message)
}
