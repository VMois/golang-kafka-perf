package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/segmentio/kafka-go"
)

const kafkaAddress string = "localhost:9093"
const topic string = "metrics"
const topicPartition int = 0
const channelSize int = 10

type MetricsData struct {
	Type     string `json:"type" binding:"required"`
	ClientID int    `json:"client_id" binding:"required"`
}

func producer(channel chan MetricsData) {
	conn, err := kafka.DialLeader(context.Background(), "tcp", kafkaAddress, topic, topicPartition)
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}

	defer conn.Close()

	for metric := range channel {
		data, err := json.Marshal(metric)
		if err != nil {
			// TODO: a better way to log errors?
			fmt.Println(err.Error())
		}

		_, err = conn.WriteMessages(
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

		c.Status(http.StatusOK)
	})

	r.Run("localhost:8080")
}
