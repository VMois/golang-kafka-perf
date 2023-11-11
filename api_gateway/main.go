package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

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

var red int = 0
var yellow int = 0
var green int = 0

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
	var err error
	var data []byte
	const retries = 3

	for metric := range channel {
		data, err = json.Marshal(metric)
		if err != nil {
			// TODO: a better way to log errors?
			fmt.Println(err.Error())
		}

		for i := 0; i < retries; i++ {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			// attempt to create topic prior to publishing the message
			err = w.WriteMessages(ctx, kafka.Message{Value: data})
			if errors.Is(err, kafka.LeaderNotAvailable) || errors.Is(err, context.DeadlineExceeded) {
				time.Sleep(time.Millisecond * 250)
				continue
			}

			if err != nil {
				fmt.Println(err.Error())
			}

			switch metric.Type {
			case "red":
				red++
			case "yellow":
				yellow++
			case "green":
				green++
			}
			break
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

	r.GET("/statistics", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"red":    red,
			"green":  green,
			"yellow": yellow,
		})
	})

	r.Run("0.0.0.0:8080")
}
