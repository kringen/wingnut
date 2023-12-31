package main

import (
	"encoding/json"
	"os"
	"time"

	rmq "github.com/kringen/message-center/rabbitmq"
)

type Status struct {
	Running bool   `json:"running"`
	Memory  string `json:"memory"`
	Cpu     string `json:"cpu"`
	Temp    string `json:"temp"`
}

func StatusCheck() {
	time.Sleep(100 * time.Second)
	// Establish messaging connection
	messageCenter := rmq.MessageCenter{}
	// Define RabbitMQ server URL.
	messageCenter.ServerUrl = os.Getenv("RABBIT_URL")
	channelName := "wingnut"
	err := messageCenter.Connect(channelName, 5, 5)
	if err != nil {
		panic(err)
	}
	defer messageCenter.Connection.Close()
	defer messageCenter.Channel.Close()
	// Create a status message
	currentStatus := Status{
		Temp: "20",
	}
	b, err := json.Marshal(currentStatus)
	if err != nil {
		panic(err)
	}

	for {
		publishMessage(&messageCenter, "health", b)
		time.Sleep(60 * time.Second)
	}
}
