package main

import (
	"log"
	"net"
	"os"

	pb "github.com/tonytheleg/grpc-go/proto/todo/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type server struct {
	d db
	pb.UnimplementedTodoServiceServer
}

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		log.Fatal("usage: server [IP_ADDR]")
	}

	addr := args[0]

	// setup a listener for the grpc server
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v\n", err)
	}

	defer func(lis net.Listener) {
		if err := lis.Close(); err != nil {
			log.Fatalf("unexpected error: %v", err)
		}
	}(lis)
	log.Printf("listening on %s\n", addr)

	creds, err := credentials.NewServerTLSFromFile("certs/server_cert.pem", "certs/server_key.pem")
	if err != nil {
		log.Fatalf("failed to create credentials %v", err)
	}
	// setup grpc server options and server
	opts := []grpc.ServerOption{
		grpc.Creds(creds),
		grpc.ChainUnaryInterceptor(unaryAuthInterceptor, unaryLogInterceptor),
		grpc.ChainStreamInterceptor(streamAuthInterceptor, streamLogInterceptor),
	}
	s := grpc.NewServer(opts...)

	pb.RegisterTodoServiceServer(s, &server{d: newDb()})

	defer s.Stop()
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v\n", err)
	}
}
