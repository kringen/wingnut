package main

import (
	"encoding/json"
	"fmt"

	rmq "github.com/kringen/message-center/rabbitmq"
)

type Configuration struct {
	Mode      string `json:"mode"`
	Objective string `json:"objective"`
}

func ConfigListen(configChan chan string, messageCenter *rmq.MessageCenter, queue string) {
	// Listen for messages
	messages, err := messageCenter.Channel.Consume(
		queue, // queue name
		"",    // consumer
		true,  // auto-ack
		false, // exclusive
		false, // no local
		false, // no wait
		nil,   // arguments
	)
	if err != nil {
		logger.Error(err.Error())
	}

	logger.Info("Successfully connected to RabbitMQ")
	logger.Info("Waiting for messages")
	for message := range messages {
		// Continue to receive messages
		logger.Info(fmt.Sprintf(" > Received message: %s\n", message.Body))
		var config Configuration
		err := json.Unmarshal(message.Body, &config)
		if err != nil {
			logger.Error(err.Error())
		}
		logger.Info(fmt.Sprintf("Mode: %s\n", config.Mode))
		logger.Info(fmt.Sprintf("Objective: %s\n", config.Objective))
	}
}
