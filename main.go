// package main
//
// import (
//
//	"bufio"
//	"encoding/json"
//	"fmt"
//	"time"
//
//	r "github.com/adam-fraga/eris/requests"
//
// )
//
// func main() {
//
//		start := time.Now()
//		msg := r.Message{
//			Role:    "user",
//			Content: "Why is the sky blue?",
//		}
//		req := r.ChatRequest{
//			Model:    "qwen3-vl:8b",
//			Stream:   true,
//			Messages: []r.Message{msg},
//		}
//		body, err := r.SendOllamaStreamRequest("http://localhost:11434/api/chat", req)
//		if err != nil {
//			panic(err)
//		}
//		defer body.Close()
//
//		scanner := bufio.NewScanner(body)
//
//		for scanner.Scan() {
//			var chunk r.Response
//			json.Unmarshal(scanner.Bytes(), &chunk)
//
//			fmt.Print(chunk.Message.Content)
//
//			if chunk.Done {
//				break
//			}
//		}
//
//		fmt.Printf("Completed in %v", time.Since(start))
//	}
package main

import "github.com/adam-fraga/eris/cmd"

func main() {
	cmd.Execute()
}
