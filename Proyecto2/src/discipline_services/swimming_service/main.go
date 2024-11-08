package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	pb "swimming_service/proto"
	"time"

	"github.com/segmentio/kafka-go"
	"google.golang.org/grpc"
)

const discipline = "swimming"

var (
	port = flag.Int("port", 50051, "The server port")
)

// Server is used to implement the gRPC server in the proto library
type server struct {
	pb.UnimplementedStudentServer
}

// Implement the GetStudent method
func (s *server) GetStudentReq(_ context.Context, in *pb.StudentRequest) (*pb.StudentResponse, error) {
	rand.Seed(time.Now().UnixNano())
	value1 := rand.Intn(2) // Random number between 0 and 1

	if value1 == 1 {
		Produce(StudentOrder{
			Student:        in.GetStudent(),
			Faculty:        in.GetFaculty(),
			Age:            int(in.GetAge()),
			Discipline:     int(in.GetDiscipline()),
			Winner:         value1,
			DisciplineName: discipline,
		}, "winners")
	}else{
		Produce(StudentOrder{
			Student:        in.GetStudent(),
			Faculty:        in.GetFaculty(),
			Age:            int(in.GetAge()),
			Discipline:     int(in.GetDiscipline()),
			Winner:         value1,
			DisciplineName: discipline,
		}, "losers")
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

type StudentOrder struct {
	Student        string `json:"student"`         // Nombre del estudiante
	Faculty        string `json:"faculty"`         // Facultad del estudiante
	Age            int    `json:"age"`             // Edad del estudiante
	Discipline     int    `json:"discipline"`      // Disciplina del estudiante
	Winner         int    `json:"winner"`          // Indicador de si es un ganador (1 o 0)
	DisciplineName string `json:"discipline_name"` // Nombre de la disciplina
}


func Produce(value StudentOrder, topicName string) {
	// to produce messages
	topic := topicName
	partition := 0

	conn, err := kafka.DialLeader(context.Background(), "tcp", "my-cluster-kafka-bootstrap:9092", topic, partition)
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}

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
