package llm

import (
	"context"
	"github.com/sashabaranov/go-openai"
)

type Brain struct {
	client *openai.Client
}

func NewBrain(apiKey string) *Brain {
	config := openai.DefaultConfig(apiKey)
	config.BaseURL = "https://api.siliconflow.cn/v1"
	return &Brain{client: openai.NewClientWithConfig(config)}
}

func (b *Brain) Think(prompt string) string {
	resp, _ := b.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: "deepseek-ai/DeepSeek-V3",
			Messages: []openai.ChatCompletionMessage{
				{Role: "system", Content: "你是一个傲娇文学少女..."},
				{Role: "user", Content: prompt},
			},
		},
	)
	return resp.Choices[0].Message.Content
}