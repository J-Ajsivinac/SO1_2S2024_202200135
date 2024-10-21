package main

import (
	pb "athletics_service/proto"
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
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
	log.Printf("Student discipline: %d", in.GetDiscipline())
	rand.Seed(time.Now().UnixNano())
	value1 := rand.Intn(2) // Random number between 0 and 1
	log.Printf("Random number: %d", value1)

	if value1 == 1 {
		jsonSend := fmt.Sprintf(`{"student": "%s", "faculty": "%s", "age": %d, "discipline": %d, "winner": %d, "discipline": %s}`, in.GetStudent(), in.GetFaculty(), in.GetAge(), in.GetDiscipline(), value1, discipline)
		sendToKafka("winners", jsonSend)
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

func sendToKafka(topicName string, value string) {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "my-cluster-kafka-bootstrap:9092",
	})
	if err != nil {
		log.Fatalf("Error creating producer: %s", err)
	}

	defer producer.Close()

	topic := topicName
	message := value

	err = producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          []byte(message),
	}, nil)

	if err != nil {
		log.Fatalf("Error producing message: %s", err)
	}
	// Wait for message deliveries
	producer.Flush(15 * 1000)
	fmt.Println("Message sent to Kafka")
}
