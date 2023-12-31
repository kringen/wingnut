package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	rmq "github.com/kringen/message-center/rabbitmq"
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	Version = "0.0.1"
	banner  = `

 _      __(_)___  ____  ____  __  __/ /_
| | /| / / / __ \/ __ \/ __ \/ / / / __/
| |/ |/ / / / / / /_/ / / / / /_/ / /_  
|__/|__/_/_/ /_/\__, /_/ /_/\__,_/\__/  v0.0.1
               /____/                   
`
)

/*
	func PrintCompletionResponse(chanChat chan string, query string) {
		completions := CreateCompletion(query)
		logger.Println(completions)
	}
*/

var logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

func createQueue(mc *rmq.MessageCenter, name string, durable bool, deleteUnused bool,
	exclusive bool, noWait bool, arguments map[string]interface{}) error {

	err := mc.CreateQueue(name, durable, deleteUnused, exclusive, noWait, arguments)
	if err != nil {
		return err
	}
	return nil
}

func publishMessage(messageCenter *rmq.MessageCenter, q string, message []byte) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var err error
	err = messageCenter.Channel.PublishWithContext(ctx,
		"",    // exchange
		q,     // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        message,
		})
	if err != nil {
		logger.Error("Failed to publish a message")
	}
	logger.Info(fmt.Sprintf(" [x] Sent %s\n", message))
}

func main() {
	fmt.Print(banner)

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

	logger.Info("Creating queue for config...")
	err = createQueue(&messageCenter, "config", false, false, false, false, nil)
	logger.Info("Creating queue for health...")
	err = createQueue(&messageCenter, "health", false, false, false, false, nil)

	// Startup service channels
	logger.Info("Starting Services...")
	//chanChat := make(chan string)
	chanConfig := make(chan string)
	go StatusCheck()
	go ConfigListen(chanConfig, &messageCenter, "config")

	// Wait for Channel to complete
	<-chanConfig

	/*
		// Wait for input
		fmt.Println("Enter command:")
		fmt.Printf(">")
		reader := bufio.NewReader(os.Stdin)
		command, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		log.Println("Command: ", command)

			go PrintCompletionResponse(chanChat, command)


			go InitiateAPI()
			logger.Println("hello")

			chanStatus <- "You da best!"
			// Wait for Channel to complete
			//<-chanChat
			//<-chanStatus

			//completions := chat.CreateCompletion()
			//logger.Println(completions)
			/* Response should be like this
			{
			  "id": "chatcmpl-xxx",
			  "object": "chat.completion",
			  "created": 1678667132,
			  "model": "gpt-3.5-turbo-0301",
			  "usage": {
			    "prompt_tokens": 13,
			    "completion_tokens": 7,
			    "total_tokens": 20
			  },
			  "choices": [
			    {
			      "message": {
			        "role": "assistant",
			        "content": "\n\nThis is a test!"
			      },
			      "finish_reason": "stop",
			      "index": 0
			    }
			  ]
			}

			for {
				// Loop to keep program running
			}
	*/

}
