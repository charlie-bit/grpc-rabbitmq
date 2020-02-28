package Client

import (
	"RabbitMQ-five-model/grpcTest/helloworld"
	"Rabbitmq-common/utils"
	"context"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"google.golang.org/grpc"
	"time"
)

/**
	Author:charlie
	Description: 我将mq客户端单独拿出来，方便其它项目调用
	Time:2020-2-24
*/

func Receive(RabbitMQAddress string) {
	//设置连接超时 提醒
	var conn *amqp.Connection
	var err error
	maxTryNum := 9
	waitSecond := 1
	for i := 1; i <= maxTryNum; i++ {
		time.Sleep(time.Duration(waitSecond) * time.Second)
		log.Infof("#%d Try to connect to RabbitMQ...", i)
		conn, err = amqp.Dial(RabbitMQAddress)
		if err == nil {
			break
		}
		log.Errorf("Failed to connect to RabbitMQ, error: %v", err)
		waitSecond *= 2
	}
	if conn == nil {
		log.Errorf("Failed to connect to RabbitMQ %d times, exit", maxTryNum)
		return
	}
	defer conn.Close()

	ch, err := conn.Channel()
	utils.FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"rpc_queue", // name
		false,       // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)
	utils.FailOnError(err, "Failed to declare a queue")

	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	utils.FailOnError(err, "Failed to set QoS")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	utils.FailOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf(" [.] fib(%s)", d.Body)
			if string(d.Body) == "" {
				log.Println("服务端发送数据为空")
				return
			}
			response := "OK"
			err = ch.Publish(
				"",        // exchange
				d.ReplyTo, 			// routing key
				false,     // mandatory
				false,     // immediate
				amqp.Publishing{
					ContentType:   "application/json",
					CorrelationId: d.CorrelationId,
					Body:          []byte(response),
				})
			utils.FailOnError(err, "Failed to publish a message")
			d.Ack(false)
		}
	}()
	log.Printf(" [*] Awaiting RPC requests")
	<-forever
}

//适用于GRPC客户端调用，响应服务端
/**
	Author:charlie
	Description:响应服务端数据
	Time:2020-2-25
*/
func GRPC_Receive(RabbitMQAddress , FunctionName string,function func(context.Context,*helloworld.HelloRequest,...grpc.CallOption)(*helloworld.HelloReply,error),
	ctx context.Context) {
	//设置连接超时 提醒
	var conn *amqp.Connection
	var err error
	maxTryNum := 9
	waitSecond := 1
	for i := 1; i <= maxTryNum; i++ {
		time.Sleep(time.Duration(waitSecond) * time.Second)
		log.Infof("#%d Try to connect to RabbitMQ...", i)
		conn, err = amqp.Dial(RabbitMQAddress)
		if err == nil {
			break
		}
		log.Errorf("Failed to connect to RabbitMQ, error: %v", err)
		waitSecond *= 2
	}
	if conn == nil {
		log.Errorf("Failed to connect to RabbitMQ %d times, exit", maxTryNum)
		return
	}

	ch, err := conn.Channel()
	utils.FailOnError(err, "Failed to open a channel")
	q, err := ch.QueueDeclare(
		FunctionName, // name
		true,       // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)
	utils.FailOnError(err, "Failed to declare a queue")

	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	utils.FailOnError(err, "Failed to set QoS")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	utils.FailOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			resp, _ := function(ctx,&helloworld.HelloRequest{Name:string(d.Body)})
			//log.Printf("server -> client : %s",string(d.Body))
			err = ch.Publish(
				"",        // exchange
				d.ReplyTo, 			// routing key
				false,     // mandatory
				false,     // immediate
				amqp.Publishing{
					ContentType:   "application/json",
					CorrelationId: d.CorrelationId,
					Body:          []byte(resp.Msg),
				})
			utils.FailOnError(err, "Failed to publish a message")
			//log.Printf("client -> server : %s",resp.Msg)
			d.Ack(false)
		}
	}()
	log.Printf(" [*] Awaiting RPC requests")
	<-forever
}
