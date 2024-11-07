package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
)

type StudentOrder struct {
	Student        string `json:"student"`         // Nombre del estudiante
	Faculty        string `json:"faculty"`         // Facultad del estudiante
	Age            int    `json:"age"`             // Edad del estudiante
	Discipline     int    `json:"discipline"`      // Disciplina del estudiante
	Winner         int    `json:"winner"`          // Indicador de si es un ganador (1 o 0)
	DisciplineName string `json:"discipline_name"` // Nombre de la disciplina
}

type StudentSave struct {
	Faculty    string `json:"faculty"`    // Facultad del estudiante
	IsWinner   bool   `json:"is_winner"`  // Indicador de si es un ganador (true o false)
	Discipline string `json:"discipline"` // Nombre de la disciplina
}

type RedisClient struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisClient() *RedisClient {
	client := redis.NewClient(&redis.Options{
		Addr:     "redis-master:6379", // Dirección de tu servidor Redis
		Password: "Wfsitpf49L",        // Contraseña
		DB:       0,                   // Base de datos de Redis
	})
	ctx := context.Background()

	return &RedisClient{
		client: client,
		ctx:    ctx,
	}
}

func (r *RedisClient) SaveStudent(student StudentSave) error {
	// Incrementar el contador de la facultad
	_, err := r.client.HIncrBy(r.ctx, "faculty:count", student.Faculty, 1).Result()
	if err != nil {
		return fmt.Errorf("failed to save student: %s", err)
	}

	// Incrementar contadores de ganadores o perdedores
	if !student.IsWinner {
		if student.Discipline != "" {
			_, err = r.client.HIncrBy(r.ctx, "discipline:losers", student.Discipline, 1).Result()
			if err != nil {
				return fmt.Errorf("failed to save student: %s", err)
			}
			_, err = r.client.HIncrBy(r.ctx, "faculty:losers", student.Faculty, 1).Result()
			if err != nil {
				return fmt.Errorf("failed to save student: %s", err)
			}
		}
	}
	return nil
}

func processEvent(event []byte, redisClient *RedisClient) {
	var student StudentOrder
	err := json.Unmarshal(event, &student)
	if err != nil {
		log.Println("failed to unmarshal event:", err)
		return
	}

	studentSave := StudentSave{
		Faculty:    student.Faculty,
		IsWinner:   student.Winner == 1,
		Discipline: student.DisciplineName,
	}

	// Guardar el estudiante en Redis
	err = redisClient.SaveStudent(studentSave)
	if err != nil {
		log.Println("failed to save student:", err)
	}
}

func main() {
	redisClient := NewRedisClient()
	defer redisClient.client.Close() 

	topic := "losers"

	// Configuración del lector de Kafka
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{"my-cluster-kafka-bootstrap:9092"},
		Topic:       topic,
		Partition:   0,
		MinBytes:    10e3, // 10KB
		MaxBytes:    10e6, // 10MB
		StartOffset: kafka.LastOffset,
		GroupID:     "my-group",
	})

	// Bucle para leer mensajes de Kafka
	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Println("failed to read message:", err)
			break
		}

		// Procesar el evento
		processEvent(m.Value, redisClient)

		// Confirmar el mensaje
		err = r.CommitMessages(context.Background(), m)
		if err != nil {
			log.Println("failed to commit message:", err)
		}
	}

	if err := r.Close(); err != nil {
		log.Fatal("failed to close reader:", err)
	}
}
