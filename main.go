package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	chat "github.com/kringen/wingnut/chat"
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

func PrintCompletionResponse(chanChat chan string, query string) {
	completions := chat.CreateCompletion(query)
	logger.Println(completions)
}

func StatusCheck(status chan string) {
	//time.Sleep(100 * time.Second)
	currentStatus := <-status
	logger.Println(currentStatus)
}

func main() {

	logger.Printf(banner, Version)

	// Startup service channels
	logger.Println("Starting Services...")
	chanChat := make(chan string)
	chanStatus := make(chan string)

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
	go StatusCheck(chanStatus)
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
	*/

	for {
		// Loop to keep program running
	}

}
