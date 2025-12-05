package requests

import (
	t "github.com/adam-fraga/eris/tools"
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

func addToolResult(req *ChatRequest, tc t.ToolCall, result string) {
	req.Messages = append(req.Messages, Message{
		Role:     "tool",
		ToolName: tc.Function.Name, // back‑reference if you want
		Content:  result,
	})
}
