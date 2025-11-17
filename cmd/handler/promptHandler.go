package handler

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"time"

	c "github.com/adam-fraga/eris/capabilities"
	"github.com/adam-fraga/eris/prompts"
	r "github.com/adam-fraga/eris/requests"
)

func RunPrompt() {
	start := time.Now()

	// Ask user for input
	fmt.Print("Ask your question: ")
	scannerInput := bufio.NewScanner(os.Stdin)
	if !scannerInput.Scan() {
		fmt.Println("No input detected")
		return
	}
	userInput := scannerInput.Text()

	userMsg := r.Message{
		Role:    "user",
		Content: userInput,
	}

	systemMsg := r.Message{
		Role:    "system",
		Content: prompts.BuildSystemPrompt(),
	}

	req := r.ChatRequest{
		Model:    "qwen3-vl:8b",
		Stream:   true,
		Messages: []r.Message{systemMsg, userMsg},
		Tools:    c.ToChatCapabilitiesInterface(),
	}

	body, err := r.SendOllamaStreamRequest("http://localhost:11434/api/chat", req)
	if err != nil {
		panic(err)
	}
	defer body.Close()

	scanner := bufio.NewScanner(body)
	for scanner.Scan() {
		var chunk r.Response
		json.Unmarshal(scanner.Bytes(), &chunk)

		// Print assistant text
		if chunk.Message.Content != "" {
			fmt.Print(chunk.Message.Content)
		}

		// Process tool calls (Ollama streaming)
		for _, tool := range chunk.Message.ToolCalls {
			argsJson, _ := json.Marshal(tool.Function.Arguments)
			result, err := c.CallCapability(tool.Function.Name, argsJson)
			if err != nil {
				fmt.Println("\n[Capability error]:", err)
			} else {
				fmt.Println("\n[Capability result]:", result)
			}
		}

		if chunk.Done {
			break
		}
	}

	fmt.Printf("\nCompleted in %v\n", time.Since(start))
}
