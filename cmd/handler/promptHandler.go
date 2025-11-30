package handler

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"time"

	c "github.com/adam-fraga/eris/config"
	r "github.com/adam-fraga/eris/requests"
)

func RunPrompt() error {
	cfg, err := c.LoadConfig()
	ctx := context.Background()

	if err != nil {
		return fmt.Errorf("error reading config file: %w", err)
	}

	model := cfg.Model
	url := cfg.Url
	// systemPrompt := cfg.SystemPrompt

	start := time.Now()

	// read the user query
	fmt.Print("Ask your question: ")
	scannerInput := bufio.NewScanner(os.Stdin)
	if !scannerInput.Scan() {
		return fmt.Errorf("failed to read user input")
	}
	userInput := scannerInput.Text()

	req := r.ChatRequest{
		Model: model,
		Messages: []r.Message{
			{Role: "system", Content: cfg.SystemPrompt},
			{Role: "user", Content: userInput},
		},
		Tools:  []r.Tool{r.GetTemperatureTool(), r.CreateFileTool()},
		Stream: true,
		Think:  true,
	}

	if err := r.SendOllamaStreaming(ctx, url, req); err != nil {
		fmt.Println("❌", err)
	}

	fmt.Printf("\nCompleted in %v\n", time.Since(start))
	return nil
}
