package main

import (
	"context"
	"flag"
	"fmt"
	pb "go_service/proto"
	"go_service/schemas"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
)

func sendData(fiberCtx *fiber.Ctx) error {
	var body schemas.Student
	if err := fiberCtx.BodyParser(&body); err != nil {
		return fiberCtx.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Set up a connection to the server.
	conn, err := grpc.NewClient(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewStudentClient(conn)

	// Create a channel to receive the response and error
	responseChan := make(chan *pb.StudentResponse)
	errorChan := make(chan error)
	go func() {

		// Contact the server and print out its response.
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		r, err := c.GetStudentReq(ctx, &pb.StudentRequest{
			Student:       body.Student,
			Age:        int32(body.Age),
			Faculty:    body.Faculty,
			Discipline: pb.Discipline(body.Discipline),
		})

		if err != nil {
			errorChan <- err
			return
		}

		responseChan <- r
	}()
	fmt.Println(len(responseChan))
	select {
	case response := <-responseChan:
		return fiberCtx.JSON(fiber.Map{
			"message": response.GetSuccess(),
		})
	case err := <-errorChan:
		return fiberCtx.Status(500).JSON(fiber.Map{
			"error": err.Error(),
			"message": "failed",
		})
	case <-time.After(5 * time.Second):
		return fiberCtx.Status(500).JSON(fiber.Map{
			"error": "timeout",
		})
	}
	
	// return fiberCtx.JSON(fiber.Map{
	// 	"message": "success",
	// })
}

func main() {
	app := fiber.New()
	app.Post("/faculty", sendData)

	err := app.Listen(":8080")
	if err != nil {
		log.Println(err)
		return
	}
}