package server

import (
	"Rabbitmq-common/utils"
	"github.com/gin-gonic/gin"
	"github.com/streadway/amqp"
	"log"
	"math/rand"
)

/**
	Author:charlie
	Description: 我将mq服务端单独拿出来，方便其它项目调用
	Time:2020-2-24
*/
func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

func randomString(l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(randInt(65, 90))
	}
	return string(bytes)
}
//适用与gin restful api 调用，将服务端作为中间件
func MqRpcServerGinMiddle(RabbitMQAddress string) gin.HandlerFunc {
	return func(context *gin.Context) {
		conn, err := amqp.Dial(RabbitMQAddress)
		utils.FailOnError(err, "Failed to connect to RabbitMQ")
		defer conn.Close()

		ch, err := conn.Channel()
		utils.FailOnError(err, "Failed to open a channel")
		defer ch.Close()

		q, err := ch.QueueDeclare(
			"",    // name
			true, // durable
			false, // delete when unused
			true,  // exclusive
			false, // noWait
			nil,   // arguments
		)
		utils.FailOnError(err, "Failed to declare a queue")

		msgs, err := ch.Consume(
			q.Name, // queue
			"",     // consumer
			false,   // auto-ack
			false,  // exclusive
			false,  // no-local
			false,  // no-wait
			nil,    // args
		)
		utils.FailOnError(err, "Failed to register a consumer")

		corrId := randomString(32)
		message := "charlie"
		err = ch.Publish(
			"",          // exchange
			"rpc_queue", // routing key
			false,       // mandatory
			false,       // immediate
			amqp.Publishing{
				ContentType:   "text/plain",
				CorrelationId: corrId,
				ReplyTo:       q.Name,
				Body:          []byte(message),
			})
		utils.FailOnError(err, "Failed to publish a message")
		log.Printf("successfully server send message : %s",message)
		for d := range msgs {
				log.Printf("client send message : %s",string(d.Body))
				resp := string(d.Body)
			if resp == "OK" {
				context.Next()
				return
				}
			}
		}
}

//适用于GRPC微服务调用，将服务请求放进消息队列，然后响应
/**
	Author:charlie
	Description:GRPC是要保证服务端客户端连通的情况下，将请求服务推送到mq，再由mq转发到服务具体接口
	Time:2020-2-25
*/
var connection *amqp.Connection
var channel *amqp.Channel
//初始化消息队列
func GRPC_SetUpMQ(RabbitAddress string) (err error) {
	if channel == nil {
		log.Println("try to connect to RabbitMQ...")
		connection, err = amqp.Dial(RabbitAddress)
		utils.FailOnError(err,"fail to connect RabbitMQ")
		channel,err = connection.Channel()
		utils.FailOnError(err,"fail to open channel")
	}
	return err
}

//将服务放进消息队列
//传入的参数是：msg：需要传入的json数据；
func GRPC_PutMSG(msg ,FunctionName string) (resp string) {
	GRPC_SetUpMQ("amqp://guest:guest@localhost:5672/")
	defer connection.Close()
	q, err := channel.QueueDeclare(
		FunctionName,    // name
		true, // durable
		false, // delete when unused
		false,  // exclusive
		false, // noWait
		nil,   // arguments
	)
	defer channel.Close()
	utils.FailOnError(err, "Failed to declare a queue")
	msgs, err := channel.Consume(
		q.Name, // queue
		"",     // consumer
		false,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	utils.FailOnError(err, "Failed to register a consumer")

	corrId := randomString(32)
	err = channel.Publish(
		"",          // exchange
		FunctionName, // routing key
		false,       // mandatory
		false,       // immediate
		amqp.Publishing{
			ContentType:   "application/json",
			CorrelationId: corrId,
			ReplyTo:       q.Name,
			Body:          []byte(msg),
		})
	utils.FailOnError(err, "Failed to publish a message")
	log.Printf("successfully server send message : %s",msg)
	for d := range msgs {
		if string(d.Body) == msg {
			return ""
		}
		resp = string(d.Body)
		return resp
	}
	return "client dont finish task,please try again soon"
}
