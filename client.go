package main

import (
	pb "bgrpcstream/proto"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"io"
	"log"
)

type Gclient struct {
	client pb.StringServicesClient
}

func Makeclient (address string) (*Gclient , bool) {
	con , err := grpc.Dial(address , grpc.WithInsecure())
	if err != nil {
		log.Println(err)
		return nil , false
	}
	client := pb.NewStringServicesClient(con)
	return &Gclient{client: client} , true
}

func (s Gclient) Dcancat (request pb.StringRequest) (*pb.StringResponse , error){
	rep , err := s.client.Concat(context.Background() , &request)
	if err != nil {
		log.Println(err)
		return nil , err
	}
	return rep , nil
}
func (s Gclient) DLotsOfserverStream (request pb.StringRequest) ([]*pb.StringResponse , error) {
	rep , err := s.client.LotsOfserverStream(context.Background() , &request)
	if err != nil {
		return nil, err
	}
	var params  = make([]*pb.StringResponse , 0)
	for {
		stream , err := rep.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil , err
		}
		params = append(params , stream)
	}
	return params , nil
}
func (s Gclient) DLostR (request ...*pb.StringRequest) (*pb.StringResponse , error) {
	client , err :=s.client.LostR(context.Background())
	if err != nil {
	return nil , err
	}
	num := len(request)
	for i := 0 ; i < num ; i ++ {
		client.Send(request[i])
	}
	rep , err := client.CloseAndRecv()
	if err != nil {
		return nil ,err
	}
	return rep , nil
}
func (s Gclient) DLostRe (request ...*pb.StringRequest) ([]*pb.StringResponse , error) {
	num := len(request)
	var params = make([]*pb.StringResponse , 0)
	client , err := s.client.LostRe(context.Background())
	if err != nil {
		return nil , err
	}
	for i := 0 ; i < num ; i ++ {
		if err := client.Send(request[i]) ; err != nil {
			return nil , err
		}
		rep , err := client.Recv()
		if err != nil {
			return nil , err
		}
		params = append(params , rep)
	}
	return params , nil
}
func main () {
	client ,ok:= Makeclient(":8083")
	if !ok {
		log.Println("dial with remote address fault")
	}
	request := &pb.StringRequest{A: "a" , B: "b"}
	rep , err := client.DLostRe(request)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(rep)
}