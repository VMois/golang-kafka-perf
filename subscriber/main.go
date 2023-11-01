package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/segmentio/kafka-go"
)

const topic string = "metrics"
const maxBatchSize int = 1e6 // 1MB

// TODO: duplicate from api gateway
type MetricsData struct {
	Type     string `json:"type" binding:"required"`
	ClientID string `json:"client_id" binding:"required"`
}

func logf(msg string, a ...interface{}) {
	fmt.Printf(msg, a...)
	fmt.Println()
}

func subscriber() {
	time.Sleep(10)
	kafkaAddress := os.Getenv("KAFKA_ADDRESS")
	if kafkaAddress == "" {
		kafkaAddress = "localhost:9093"
	}
	fmt.Println("Trying connect to Kafka at", kafkaAddress)

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{kafkaAddress},
		Topic:       topic,
		MaxBytes:    maxBatchSize,
		MaxAttempts: 10,
		GroupID:     "consumer-group-id",
		Logger:      kafka.LoggerFunc(logf),
		ErrorLogger: kafka.LoggerFunc(logf),
	})

	// TODO: will not work with SIGINT
	defer r.Close()

	fmt.Println("Consumer ready for accepting data")

	var metric MetricsData
	seen := make(map[string]bool)

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			fmt.Println("Couldn't read message. Msg: ", err.Error())
			time.Sleep(1)
			continue
		}
		err = json.Unmarshal(m.Value, &metric)
		if err != nil {
			fmt.Println("Couldn't unmarshall data. Msg: ", err.Error())
			time.Sleep(1)
			continue
		}

		if seen[metric.ClientID] {
			fmt.Println("Found a duplicate of ", metric.ClientID)
		} else {
			seen[metric.ClientID] = true
		}

		fmt.Println("Received metric: ", metric.Type, metric.ClientID)

		switch metric.Type {
		case "red":
			red++
		case "yellow":
			yellow++
		case "green":
			green++
		}
	}
}

// TODO: Global vars, hell yeah!
var red int = 0
var yellow int = 0
var green int = 0

func main() {
	r := gin.Default()
	go subscriber()

	r.GET("/statistics", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"red":    red,
			"yellow": yellow,
			"green":  green,
		})
	})

	r.GET("/health", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	r.POST("/reset", func(c *gin.Context) {
		red = 0
		green = 0
		yellow = 0
		c.Status(http.StatusOK)
	})

	r.Run("0.0.0.0:8080")
}
