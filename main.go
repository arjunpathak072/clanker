package main

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"google.golang.org/genai"
)

func main() {
	godotenv.Load()
	ctx := context.Background()
	client, err := genai.NewClient(ctx, nil)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return
	}

	scanner := bufio.NewScanner(os.Stdin)
	getUserMessage := func() (string, bool) {
		if !scanner.Scan() {
			return "", false
		}
		return scanner.Text(), true
	}

	agent := NewAgent(client, getUserMessage)
	err = agent.Run(context.TODO())
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
	}
}

func NewAgent(client *genai.Client, getUserMessage func() (string, bool)) *Agent {
	return &Agent{
		client:         client,
		getUserMessage: getUserMessage,
	}
}

type Agent struct {
	client         *genai.Client
	getUserMessage func() (string, bool)
}

func (a *Agent) Run(ctx context.Context) error {
	conversation := []*genai.Content{}

	fmt.Println("Chat with Gemini (use 'ctrl-c' to quit)")

	for {
		fmt.Print("\u001b[94mYou\u001b[0m: ")
		userInput, ok := a.getUserMessage()
		if !ok {
			break
		}

		userMessage := &genai.Content{
			Role:  "user",
			Parts: []*genai.Part{{Text: userInput}},
		}
		conversation = append(conversation, userMessage)

		message, err := a.runInference(ctx, conversation)
		if err != nil {
			return err
		}
		conversation = append(conversation, message.Candidates[0].Content)

		fmt.Printf("\u001b[93mGemini\u001b[0m: %s\n", message.Text())
	}

	return nil
}

func (a *Agent) runInference(ctx context.Context, conversation []*genai.Content) (*genai.GenerateContentResponse, error) {
	message, err := a.client.Models.GenerateContent(ctx, "gemini-3-flash-preview", conversation, nil)
	return message, err
}
