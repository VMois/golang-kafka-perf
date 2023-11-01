package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/segmentio/kafka-go"
)

const topic string = "metrics"
const topicPartition int = 0
const channelSize int = 100

type MetricsData struct {
	Type     string `json:"type" binding:"required"`
	ClientID string `json:"client_id" binding:"required"`
}

func producer(channel chan MetricsData) {
	kafkaAddress := os.Getenv("KAFKA_ADDRESS")
	if kafkaAddress == "" {
		kafkaAddress = "localhost:9093"
	}

	fmt.Println("Trying connect to Kafka at", kafkaAddress)

	w := &kafka.Writer{
		Addr:                   kafka.TCP(kafkaAddress),
		Topic:                  topic,
		Balancer:               &kafka.LeastBytes{},
		AllowAutoTopicCreation: true,
	}

	defer w.Close()

	for metric := range channel {
		data, err := json.Marshal(metric)
		if err != nil {
			// TODO: a better way to log errors?
			fmt.Println(err.Error())
		}

		err = w.WriteMessages(context.Background(),
			kafka.Message{Value: data},
		)
		if err != nil {
			// TODO: What to do if error? Should we repeat?
			fmt.Println("Failed to send message to Kafka: ", err.Error())
		}
	}
}

func main() {
	r := gin.Default()
	var channel chan MetricsData = make(chan MetricsData, channelSize)
	go producer(channel)

	defer close(channel)

	r.POST("/metrics", func(c *gin.Context) {
		var metric MetricsData

		err := c.BindJSON(&metric)
		if err != nil {
			// TODO: not great way, but will work for now
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": err.Error(),
			})
			return
		}

		channel <- metric

		c.JSON(http.StatusOK, gin.H{
			"type": metric.Type,
		})
	})

	r.GET("/health", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	r.DELETE("/purge", func(c *gin.Context) {
		kafkaAddress := os.Getenv("KAFKA_ADDRESS")
		if kafkaAddress == "" {
			kafkaAddress = "localhost:9093"
		}
		conn, err := kafka.Dial("tcp", kafkaAddress)
		if err != nil {
			fmt.Println(err.Error())
		}
		defer conn.Close()

		err = conn.DeleteTopics(topic)
		if err != nil {
			fmt.Println(err.Error())
		}
	})

	r.Run("0.0.0.0:8080")
}
