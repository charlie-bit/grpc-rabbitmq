package main

import (
	"Rabbitmq-common/server"
	"log"
)
func main()  {
	name := "world"
	resp := server.GRPC_PutMSG(name,"SayHello")
	if resp == ""{
		log.Println("grpc_server dont response")
	}else {
		log.Printf("grpc_server -> grpc_client : %s", resp)
	}
}
