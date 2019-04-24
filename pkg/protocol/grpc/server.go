package grpc

import (
	"context"
	v1 "go.xiaosongfu.com/go-grpc-http-rest-microservice-tutorial/pkg/api/v1"
	"google.golang.org/grpc"
	"log"
	"net"
)

func RunServer(ctx context.Context, v1API v1.ToDoServiceServer, port string) error {
	listen, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	server := grpc.NewServer()
	v1.RegisterToDoServiceServer(server, v1API)

	log.Println("starting gRPC server...")
	return server.Serve(listen)
}