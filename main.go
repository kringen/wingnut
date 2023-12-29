package main

import (
	"encoding/json"
	"log"
	"os"
	"time"

	rmq "github.com/kringen/message-center/rabbitmq"
)

var logger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)

const (
	Version = "0.0.1"
	banner  = `

 _      __(_)___  ____  ____  __  __/ /_
| | /| / / / __ \/ __ \/ __ \/ / / / __/
| |/ |/ / / / / / /_/ / / / / /_/ / /_  
|__/|__/_/_/ /_/\__, /_/ /_/\__,_/\__/  v%s
               /____/                   
`
)

/*
	func PrintCompletionResponse(chanChat chan string, query string) {
		completions := CreateCompletion(query)
		logger.Println(completions)
	}
*/

type Configuration struct {
	Mode      string `json:"mode"`
	Objective string `json:"objective"`
}

func StatusCheck(status chan string) {
	time.Sleep(100 * time.Second)
	currentStatus := <-status
	logger.Println(currentStatus)
}

func SetConfig(configChan chan string, messageCenter *rmq.MessageCenter, queue string) {
	// Subscribing to QueueService1 for getting messages.
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
		log.Println(err)
	}

	// Build a welcome message.
	log.Println("Successfully connected to RabbitMQ")
	log.Println("Waiting for messages")
	for message := range messages {
		// For example, show received message in a console.
		log.Printf(" > Received message: %s\n", message.Body)
		var config Configuration
		err := json.Unmarshal(message.Body, &config)
		if err != nil {
			log.Println(err)
		}
		log.Printf("Mode: %s\n", config.Mode)
		log.Printf("Objective: %s\n", config.Objective)
	}
}

func createQueue(mc *rmq.MessageCenter, name string, durable bool, deleteUnused bool,
	exclusive bool, noWait bool, arguments map[string]interface{}) error {

	err := mc.CreateQueue(name, durable, deleteUnused, exclusive, noWait, arguments)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	logger.Printf(banner, Version)

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

	logger.Printf("Creating queue for config...")
	err = createQueue(&messageCenter, "config", false, false, false, false, nil)

	// Startup service channels
	logger.Println("Starting Services...")
	//chanChat := make(chan string)
	chanStatus := make(chan string)
	chanConfig := make(chan string)
	go StatusCheck(chanStatus)
	go SetConfig(chanConfig, &messageCenter, "config")

	// Wait for Channel to complete
	<-chanStatus
	<-chanConfig

	chanStatus <- "You da best!"
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
