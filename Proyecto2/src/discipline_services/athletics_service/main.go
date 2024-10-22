package main

import (
	pb "athletics_service/proto"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"time"

	"github.com/segmentio/kafka-go"
	"google.golang.org/grpc"
)

const discipline = "athletics"

var (
	port = flag.Int("port", 50051, "The server port")
)

// Server is used to implement the gRPC server in the proto library
type server struct {
	pb.UnimplementedStudentServer
}

// Implement the GetStudent method
func (s *server) GetStudentReq(_ context.Context, in *pb.StudentRequest) (*pb.StudentResponse, error) {
	log.Printf("Received: %v", in)
	log.Printf("Student name: %s", in.GetStudent())
	log.Printf("Student faculty: %s", in.GetFaculty())
	log.Printf("Student age: %d", in.GetAge())
	log.Printf("Student discipline: %s", in.GetDiscipline())

	rand.Seed(time.Now().UnixNano())
	value1 := rand.Intn(2) // Random number between 0 and 1
	log.Printf("Random number: %d", value1)

	if value1 == 1 {
		// Crear un mapa para el mensaje JSON
		// Enviar el JSON a Kafka
		// sendToKafka("winners", string(jsonSend))
		Produce(StudentOrder{
			Student:        in.GetStudent(),
			Faculty:        in.GetFaculty(),
			Age:            int(in.GetAge()),
			Discipline:     int(in.GetDiscipline()),
			Winner:         value1,
			DisciplineName: discipline,
		}, "winners")
	}

	return &pb.StudentResponse{
		Success: true,
	}, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterStudentServer(s, &server{})
	log.Printf("Server started on port %d", *port)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

// func sendToKafka(topicName string, value string) {
// 	// Configuración del productor para Strimzi Kafka
// 	producer, err := kafka.NewProducer(&kafka.ConfigMap{
// 		// Usa el FQDN completo incluyendo el namespace correcto
// 		"bootstrap.servers": "my-cluster-kafka-bootstrap.so1-proyecto2.svc.cluster.local:9092",
// 		"client.id":         "athletics-service",
// 		"acks":              "all",
// 		"retries":           3,
// 		"retry.backoff.ms":  1000,
// 		"socket.timeout.ms": 30000,
// 		"message.timeout.ms": 30000,
// 		"debug":             "broker,topic,msg",
// 	})
// 	if err != nil {
// 		log.Printf("Error creating producer: %s", err)
// 		return
// 	}
// 	defer producer.Close()

// 	// Crear el topic si no existe (opcional)
// 	adminClient, err := kafka.NewAdminClient(&kafka.ConfigMap{
// 		"bootstrap.servers": "my-cluster-kafka-bootstrap.so1-proyecto2.svc.cluster.local:9092",
// 	})
// 	if err == nil {
// 		// Intenta crear el topic
// 		topicSpecs := []kafka.TopicSpecification{{
// 			Topic:             topicName,
// 			NumPartitions:     1,
// 			ReplicationFactor: 1,
// 		}}
// 		_, err := adminClient.CreateTopics(context.Background(), topicSpecs)
// 		if err != nil {
// 			log.Printf("Warning: Could not create topic: %v", err)
// 		}
// 		adminClient.Close()
// 	}

// 	deliveryChan := make(chan kafka.Event)
// 	topic := topicName
// 	message := value

// 	err = producer.Produce(&kafka.Message{
// 		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
// 		Value:          []byte(message),
// 		Headers: []kafka.Header{
// 			{
// 				Key:   "source",
// 				Value: []byte("athletics-service"),
// 			},
// 		},
// 	}, deliveryChan)

// 	if err != nil {
// 		log.Printf("Error producing message: %s", err)
// 		return
// 	}

// 	// Manejo de eventos de entrega
// 	go func() {
// 		for e := range deliveryChan {
// 			switch ev := e.(type) {
// 			case *kafka.Message:
// 				if ev.TopicPartition.Error != nil {
// 					log.Printf("Failed to deliver message: %v\n", ev.TopicPartition.Error)
// 				} else {
// 					log.Printf("Successfully delivered message to topic %s [%d] at offset %v\n",
// 						*ev.TopicPartition.Topic, ev.TopicPartition.Partition, ev.TopicPartition.Offset)
// 				}
// 			}
// 		}
// 	}()

// 	remaining := producer.Flush(15 * 1000)
// 	if remaining > 0 {
// 		log.Printf("Failed to flush all messages. %d messages remaining\n", remaining)
// 	}

// 	close(deliveryChan)
// }

type StudentOrder struct {
	Student        string `json:"student"`         // Nombre del estudiante
	Faculty        string `json:"faculty"`         // Facultad del estudiante
	Age            int    `json:"age"`             // Edad del estudiante
	Discipline     int    `json:"discipline"`      // Disciplina del estudiante
	Winner         int    `json:"winner"`          // Indicador de si es un ganador (1 o 0)
	DisciplineName string `json:"discipline_name"` // Nombre de la disciplina
}

// usando segmentio/kafka-go

func Produce(value StudentOrder, topicName string) {
	// to produce messages
	topic := topicName
	partition := 0

	conn, err := kafka.DialLeader(context.Background(), "tcp", "my-cluster-kafka-bootstrap:9092", topic, partition)
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}

	// Convert the data struct into a byte slice
	valueBytes, err := json.Marshal(value)
	if err != nil {
		log.Fatalf("Failed to marshal value: %v", err)
	}
	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	_, err = conn.WriteMessages(
		kafka.Message{Value: valueBytes},
	)
	if err != nil {
		log.Fatal("failed to write messages:", err)
	}

	if err := conn.Close(); err != nil {
		log.Fatal("failed to close writer:", err)
	}

	log.Println("Message sent")
}
