package main

import (
	"RabbitMQ-five-model/grpcTest/helloworld"
	"Rabbitmq-common/Client"
	"context"
	"flag"
	"google.golang.org/grpc"
	"log"
)

var (
	flagRabbitMQ  = flag.String("rabbitmq", "amqp://guest:guest@localhost:5672/", "RabbitMQ Address")
)

func main()  {
	conn,err := grpc.Dial("127.0.0.1:8000",grpc.WithInsecure())
	if err != nil {
		log.Println("建立连接失败")
	}
	defer conn.Close()
	client := helloworld.NewHelloWorldClient(conn)
	Client.GRPC_Receive(*flagRabbitMQ,"SayHello",client.SayHello,context.Background())
}
