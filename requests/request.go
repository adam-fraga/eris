package requests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	t "github.com/adam-fraga/eris/tools"
	"github.com/fatih/color"
)

// The request / response schema – exactly what Ollama api expects
type ChatRequest struct {
	Model    string    `json:"model"`
	Prompt   string    `json:"prompt,omitempty"`
	System   string    `json:"system,omitempty"`
	Options  *Options  `json:"options,omitempty"`
	Messages []Message `json:"messages,omitempty"`
	Tools    []t.Tool  `json:"tools,omitempty"`
	Stream   bool      `json:"stream,omitempty"`
	Think    bool      `json:"think,omitempty"`
	Params   any       `json:"params,omitempty"` // arbitrary extra data
}

type Options struct {
	Temperature   float64 `json:"temperature,omitempty"`
	Seed          int     `json:"seed,omitempty"`
	TopK          int     `json:"top_k,omitempty"`
	TopP          float64 `json:"top_p,omitempty"`
	MinLength     int     `json:"min_length,omitempty"`
	MaxTokens     int     `json:"max_tokens,omitempty"`
	RepeatPenalty float64 `json:"repeat_penalty,omitempty"`
	RepeatLastN   int     `json:"repeat_last_n,omitempty"`
}

type Message struct {
	Role      string       `json:"role"`
	Content   string       `json:"content,omitempty"`
	Thinking  string       `json:"thinking,omitempty"`
	ToolCalls []t.ToolCall `json:"tool_calls,omitempty"`
	ToolName  string       `json:"tool_name,omitempty"` // <- added
}

// Streaming engine
func SendOllamaStreaming(ctx context.Context, url string, req ChatRequest) error {
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
				city := tc.Function.Arguments["city"].(string)
				res := t.GetTemperature(city)
				addToolResult(&req, tc, res)

			case "create_file":
				path := tc.Function.Arguments["path"].(string)
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

func addToolResult(req *ChatRequest, tc t.ToolCall, result string) {
	req.Messages = append(req.Messages, Message{
		Role:     "tool",
		ToolName: tc.Function.Name, // back‑reference if you want
		Content:  result,
	})
}
