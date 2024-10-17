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

	"google.golang.org/grpc"
)

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
