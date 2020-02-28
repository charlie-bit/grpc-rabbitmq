package main

import (
	"RabbitMQ-five-model/grpcTest/controller"
	"Rabbitmq-common/server"
	"flag"
	"github.com/gin-gonic/gin"
)

var (
	flagAddr = flag.String("addr",":8000","router listen address")
	flagRabbitMQ  = flag.String("rabbitmq", "amqp://guest:guest@localhost:5672/", "RabbitMQ Address")
)

/**
	Author:charlie
	Description: PostMan send request to router
	Time:2020-2-24
 */

func main()  {
	router := gin.New()
	v1 := router.Group("/api/v1")
	v1.Use(server.MqRpcServerGinMiddle(*flagRabbitMQ))
	{
		v1.POST("/hello",controller.Hello)
	}
	router.Run(":8000")
}
