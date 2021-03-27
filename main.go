package main

import (
	"bgrpcstream/handler"
	pb "bgrpcstream/proto"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main () {
	lis , err := net.Listen("tcp",":8083")
	if err != nil {
		log.Println(err)
	}
	grpcserver := grpc.NewServer()
	services := new(handler.StringServiceStream)
	pb.RegisterStringServicesServer(grpcserver , services)
	grpcserver.Serve(lis)
}