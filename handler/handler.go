package handler

import (
	pb "bgrpcstream/proto"
	"context"
	"io"
	"log"
	"strconv"
)

type StringServiceStream struct {}

func (s StringServiceStream) Concat (_ context.Context,request  *pb.StringRequest) (*pb.StringResponse, error) {
	rep := &pb.StringResponse{Msg: request.B+request.A}
	return rep , nil
}

func (s StringServiceStream) LotsOfserverStream (request *pb.StringRequest, qs pb.StringServices_LotsOfserverStreamServer) error {

	for i := 0 ; i < 10 ; i ++ {
		rep := &pb.StringResponse{Msg: request.A + request.B+strconv.Itoa(i)}
		qs.Send(rep)
	}
	return nil
}
func (s StringServiceStream) LostR(client pb.StringServices_LostRServer) error {
	var params []string
	for {
		stream , err := client.Recv()
		if err == io.EOF {
			client.SendAndClose(&pb.StringResponse{Msg: "end"})
			return nil
		}
		if err != nil {
			log.Println(err)
			return err
		}
		params = append(params , stream.B , stream.A)
	}
}
func (s StringServiceStream) LostRe (request pb.StringServices_LostReServer) error {
	for {
		stream , err := request.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println(err)
			return err
		}
		request.Send(&pb.StringResponse{Msg: stream.B+stream.A})
	}
	return nil
}
