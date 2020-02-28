package main

import (
	"RabbitMQ-five-model/grpcTest/controller"
	pb "RabbitMQ-five-model/grpcTest/helloworld"
	"context"
	"google.golang.org/grpc"
	"log"
	"net"
)

type ServerStruct struct {}

func (ss *ServerStruct)SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply,error)  {
	msg := in.Name
	return &pb.HelloReply{Msg:controller.SayHello(msg)},nil
}

func main()  {
	lis,err := net.Listen("tcp",":8000")
	if err != nil {
		log.Printf("fail to listen :%v",err)
	}
	s := grpc.NewServer()
	pb.RegisterHelloWorldServer(s, &ServerStruct{})

	log.Println("server run successfully...")

	if err := s.Serve(lis) ;err !=nil{
		log.Println(err)
	}
}
