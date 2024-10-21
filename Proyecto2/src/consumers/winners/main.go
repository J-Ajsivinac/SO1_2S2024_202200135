package main

import (
	"fmt"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func main() {
	kafkaBroker := "my-cluster-kafka-bootstrap:9092"
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": kafkaBroker,
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		log.Fatalf("Error creating consumer: %s", err)
	}

	defer consumer.Close()

	topic := "winners"
	err = consumer.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		log.Fatalf("Error subscribing to topic: %s", err)
	}

	fmt.Printf("Esperando mensajes en el tópico %s...\n", topic)
	for {
		msg, err := consumer.ReadMessage(-1)
		if err == nil {
			fmt.Printf("Mensaje recibido: %s\n", string(msg.Value))
		} else {
			fmt.Printf("Error leyendo mensaje: %v (%v)\n", err, msg)
		}
	}
}


