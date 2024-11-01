package main

import (
	"context"
	"flag"
	pb "go_service/proto"
	"go_service/schemas"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	athleticsAddr = flag.String("athleticsAddr", "athletics-service:50051", "athletics service address")
	swimmingAddr  = flag.String("swimmingAddr", "swimming-service:50051", "swimming service address")
	boxingAddr    = flag.String("boxingAddr", "boxing-service:50051", "boxing service address")
)

func sendDataToDiscipline(serviceAddr string, request *pb.StudentRequest, resultChan chan<- *pb.StudentResponse, errorChan chan<- error) {
	conn, err := grpc.NewClient(serviceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		errorChan <- err
		return
	}
	defer conn.Close()

	client := pb.NewStudentClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	response, err := client.GetStudentReq(ctx, request)
	if err != nil {
		errorChan <- err
		return
	}

	resultChan <- response
}

func sendData(fiberCtx *fiber.Ctx) error {
	var body schemas.Student
	if err := fiberCtx.BodyParser(&body); err != nil {
		return fiberCtx.Status(400).JSON(fiber.Map{
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
		return fiberCtx.Status(400).JSON(fiber.Map{
			"error": "Invalid discipline",
		})
	}

	// Crear los canales para manejar la respuesta y errores
	resultChan := make(chan *pb.StudentResponse)
	errorChan := make(chan error)

	// Crear una solicitud StudentRequest
	request := &pb.StudentRequest{
		Student:    body.Student,
		Faculty:    body.Faculty,
		Age:        int32(body.Age),
		Discipline: pb.Discipline(body.Discipline),
	}

	// Iniciar una goroutine para enviar la solicitud gRPC
	go sendDataToDiscipline(serviceAddr, request, resultChan, errorChan)

	select {
	case response := <-resultChan:
		return fiberCtx.Status(200).JSON(fiber.Map{
			"success": response.Success,
		})
	case err := <-errorChan:
		return fiberCtx.Status(500).JSON(fiber.Map{
			"error":   "Cannot get student",
			"message": err.Error(),
		})
	case <-time.After(time.Second):
		return fiberCtx.Status(500).JSON(fiber.Map{
			"error": "Request timeout",
		})
	}
}

func main() {
	flag.Parse()
	app := fiber.New()
	app.Post("/agronomy", sendData)

	if err := app.Listen(":8080"); err != nil {
		log.Println(err)
		return
	}
}
