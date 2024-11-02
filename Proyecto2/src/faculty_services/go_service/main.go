package main

import (
	"context"
	"flag"
	"fmt"
	pb "go_service/proto"
	"go_service/schemas"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

var (
	athleticsAddr = flag.String("athleticsAddr", "athletics-service:50051", "athletics service address")
	swimmingAddr  = flag.String("swimmingAddr", "swimming-service:50051", "swimming service address")
	boxingAddr    = flag.String("boxingAddr", "boxing-service:50051", "boxing service address")
)

type ConnectionPool struct {
	connections map[string]*grpc.ClientConn
	mu         sync.RWMutex
}

var pool = &ConnectionPool{
	connections: make(map[string]*grpc.ClientConn),
}

func (p *ConnectionPool) GetConnection(serviceAddr string) (*grpc.ClientConn, error) {
	p.mu.RLock()
	conn, exists := p.connections[serviceAddr]
	p.mu.RUnlock()

	if exists {
		return conn, nil
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	// Double-check after acquiring write lock
	if conn, exists := p.connections[serviceAddr]; exists {
		return conn, nil
	}

	kacp := keepalive.ClientParameters{
		Time:                10 * time.Second,
		Timeout:            3 * time.Second,
		PermitWithoutStream: true,
	}

	var err error
	// Eliminamos la declaración redundante de conn
	for retries := 3; retries > 0; retries-- {
		conn, err = grpc.NewClient(
			serviceAddr,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithKeepaliveParams(kacp),
			grpc.WithBlock(),
		)

		if err == nil {
			p.connections[serviceAddr] = conn
			return conn, nil
		}

		if retries > 1 {
			time.Sleep(time.Second)
			log.Printf("Retrying connection to %s, attempts remaining: %d", serviceAddr, retries-1)
		}
	}

	return nil, fmt.Errorf("failed to connect after 3 retries: %v", err)
}

func sendDataToDiscipline(ctx context.Context, serviceAddr string, request *pb.StudentRequest) (*pb.StudentResponse, error) {
	conn, err := pool.GetConnection(serviceAddr)
	if err != nil {
		return nil, fmt.Errorf("connection error: %v", err)
	}

	client := pb.NewStudentClient(conn)
	
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	response, err := client.GetStudentReq(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("RPC error: %v", err)
	}

	return response, nil
}

func sendData(fiberCtx *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(fiberCtx.Context(), 10*time.Second)
	defer cancel()

	var body schemas.Student
	if err := fiberCtx.BodyParser(&body); err != nil {
		return fiberCtx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Cannot parse body",
			"message": err.Error(),
		})
	}

	var serviceAddr string
	switch body.Discipline {
	case 1:
		serviceAddr = *swimmingAddr
	case 2:
		serviceAddr = *athleticsAddr
	case 3:
		serviceAddr = *boxingAddr
	default:
		return fiberCtx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid discipline",
		})
	}

	request := &pb.StudentRequest{
		Student:    body.Student,
		Faculty:    body.Faculty,
		Age:        int32(body.Age),
		Discipline: pb.Discipline(body.Discipline),
	}

	response, err := sendDataToDiscipline(ctx, serviceAddr, request)
	if err != nil {
		log.Printf("Error processing request: %v", err)
		return fiberCtx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Cannot process request",
			"message": err.Error(),
		})
	}

	return fiberCtx.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": response.Success,
	})
}

func main() {
	flag.Parse()

	app := fiber.New(fiber.Config{
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  30 * time.Second,
	})

	app.Use(recover.New())
	app.Post("/agronomy", sendData)

	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Println("Starting server on :8080...")
		if err := app.Listen(":8080"); err != nil {
			log.Printf("Server error: %v", err)
		}
	}()

	<-shutdownChan
	log.Println("Shutting down server...")

	for _, conn := range pool.connections {
		if err := conn.Close(); err != nil {
			log.Printf("Error closing gRPC connection: %v", err)
		}
	}

	if err := app.Shutdown(); err != nil {
		log.Printf("Error during server shutdown: %v", err)
	}

	log.Println("Server shutdown complete")
}