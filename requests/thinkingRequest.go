package requests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	t "github.com/adam-fraga/eris/tools"
	"github.com/fatih/color"
)

func SendOllamaThinkRequest(ctx context.Context, url string, req ChatRequest) error {
	// Helper that actually does the HTTP POST + streaming read
	doRequest := func(r ChatRequest) (*http.Response, error) {
		payload, _ := json.Marshal(r)
		return http.Post(url, "application/json", bytes.NewReader(payload))
	}

	// --- Outer loop: keep feeding new context until no tool calls remain
	for {
		resp, err := doRequest(req)
		if err != nil {
			return fmt.Errorf("post: %w", err)
		}
		if resp.StatusCode != http.StatusOK {
			b, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			return fmt.Errorf("non‑OK status %d: %s", resp.StatusCode, string(b))
		}
		defer resp.Body.Close()

		// Read the streaming chunks, accumulating partial fields
		var thinking, content string
		var toolCalls []t.ToolCall
		dec := json.NewDecoder(resp.Body)

		for {
			var chunk struct {
				Message Message `json:"message"`
			}
			if err := dec.Decode(&chunk); err != nil {
				if err == io.EOF {
					break
				}
				return fmt.Errorf("decode: %w", err)
			}

			if chunk.Message.Thinking != "" {

				thinking += chunk.Message.Thinking
				color.RGB(111, 122, 137).Print(chunk.Message.Thinking)
			}

			if chunk.Message.Content != "" {
				content += chunk.Message.Content
				fmt.Print(chunk.Message.Content)
			}
			if len(chunk.Message.ToolCalls) > 0 {
				toolCalls = append(toolCalls, chunk.Message.ToolCalls...)
			}
		}
		fmt.Println() // flush line after stream ends

		// Push assistant chunk to the context
		if thinking != "" || content != "" || len(toolCalls) > 0 {
			req.Messages = append(req.Messages, Message{
				Role:      "assistant",
				Thinking:  thinking,
				Content:   content,
				ToolCalls: toolCalls,
			})
		}

		// If no tool calls – conversation finished
		if len(toolCalls) == 0 {
			return nil
		}

		//  Execute each tool call and feed the result back as a new message
		for _, tc := range toolCalls {
			switch tc.Function.Name {
			case "get_temperature":
				city := strings.ToLower(tc.Function.Arguments["city"].(string))
				temp := t.GetTemperature(city)
				content := fmt.Sprintf("Current temperature in %s: %s", city, temp)
				addToolResult(&req, tc, content)

			case "create_file":
				path := tc.Function.Arguments["path"].(string)
				pathArg, ok := tc.Function.Arguments["path"].(string)
				if !ok || pathArg == "" {
					// fallback to cwd
					pathArg = "output.txt"
				}
				content := tc.Function.Arguments["content"].(string)
				res := t.CreateFile(path, content)
				addToolResult(&req, tc, res)

			default:
				addToolResult(&req, tc, "unknown function")
			}
		}

		// Build the *next* request with the updated context
		req = ChatRequest{
			Model:    req.Model,
			Prompt:   "", // empty because we already have the context
			Messages: req.Messages,
			Tools:    []t.Tool{t.GetTemperatureTool(), t.CreateFileTool()},
			Stream:   true,
			Think:    true,
		}

		// The outer loop will now send the new request again
	}
}
