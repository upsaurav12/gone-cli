package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type GroqChatRequest struct {
	Model    string        `json:"model"`
	Messages []ChatMessage `json:"messages"`
}

func CallingApiToGroq(prompt string) (string, error) {

	_ = godotenv.Load()

	url := os.Getenv("GROQ_API_URL")
	key := os.Getenv("GROQ_API_KEY")

	if url == "" || key == "" {
		return "", fmt.Errorf("missing GROQ_API_URL or GROQ_API_KEY")
	}

	body := GroqChatRequest{
		Model: "llama-3.1-8b-instant",
		Messages: []ChatMessage{
			{Role: "user", Content: prompt},
		},
	}

	jsonBody, _ := json.Marshal(body)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+key)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// READ RAW BODY FIRST
	rawBody, _ := io.ReadAll(resp.Body)
	// fmt.Println("RAW RESPONSE:", string(rawBody))

	// Decode into struct
	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.Unmarshal(rawBody, &result); err != nil {
		return "", err
	}

	if len(result.Choices) == 0 {
		return "", fmt.Errorf("no choices returned → request was invalid")
	}

	return result.Choices[0].Message.Content, nil
}
