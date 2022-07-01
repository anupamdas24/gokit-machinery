package main

import (
	"log"
	"runtime"

	"github.com/golang/protobuf/proto"
	nats "github.com/nats-io/nats.go"

	pb "github.com/shijuvar/gokit/src/grpc-nats/order"
)

const (
	queue   = "Order.OrdersCreatedQueue"
	subject = "Order.OrderCreated"
)

func main() {

	// Create server connection
	natsConnection, _ := nats.Connect(nats.DefaultURL)
	log.Println("Connected to " + nats.DefaultURL)
	// Subscribe to subject
	natsConnection.QueueSubscribe(subject, queue, func(msg *nats.Msg) {
		eventStore := pb.EventStore{}
		err := proto.Unmarshal(msg.Data, &eventStore)
		if err == nil {
			// Handle the message
			log.Printf("Subscribed message in Worker 1: %+v\n", eventStore)
		}
	})

	// Keep the connection alive
	runtime.Goexit()
}
