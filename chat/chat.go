package chat

import (
	"os"
)

func CreateCompletion(query string) CreateCompletionsResponse {
	apiKey := os.Getenv("OPENAI_KEY")
	organization := os.Getenv("OPENAI_ORG")

	client := NewClient(apiKey, organization)

	r := CreateCompletionsRequest{
		Model: "gpt-3.5-turbo",
		Messages: []Message{
			{
				Role:    "user",
				Content: query,
			},
		},
		Temperature: 0.7,
	}

	completions, err := client.CreateCompletions(r)
	if err != nil {
		panic(err)
	}

	return completions
}
