package handler

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"time"

	c "github.com/adam-fraga/eris/config"
	r "github.com/adam-fraga/eris/requests"
	t "github.com/adam-fraga/eris/tools"
)

func RunCodingPrompt() error {
	cfg, err := c.LoadConfig()
	ctx := context.Background()

	if err != nil {
		return fmt.Errorf("error reading config file: %w", err)
	}

	model := cfg.CodingModel
	url := cfg.Url
	systemPrompt := cfg.SystemPrompt

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
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userInput},
		},
		Tools:  []t.Tool{t.GetTemperatureTool(), t.CreateFileTool()},
		Stream: true,
		Think:  false,
	}

	if err := r.SendOllamaCodeRequest(ctx, url, req); err != nil {
		fmt.Println("❌", err)
	}

	fmt.Printf("\nCompleted in %v\n", time.Since(start))
	return nil
}
