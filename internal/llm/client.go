package llm

import (
	"context"
	"github.com/sashabaranov/go-openai"
)

type Brain struct {
	client  *openai.Client
	// 内存中的短期记忆：key 是用户QQ号，value 是对话历史
	history map[int64][]openai.ChatCompletionMessage
}

func NewBrain(apiKey string) *Brain {
	config := openai.DefaultConfig(apiKey)
	config.BaseURL = "https://api.siliconflow.cn/v1"
	return &Brain{
		client:  openai.NewClientWithConfig(config),
		history: make(map[int64][]openai.ChatCompletionMessage),
	}
}

func (b *Brain) Think(userID int64, prompt string) string {
	// 1. 获取该用户的历史记录，如果没有则初始化人设
	if _, ok := b.history[userID]; !ok {
		b.history[userID] = []openai.ChatCompletionMessage{
			{Role: "system", Content: "你名叫强氧化的浓硫酸，是一个傲娇文学少女。你的回复应该简洁中带有辛辣，但是充满温柔。"},
		}
	}

	// 2. 将用户的新问题加入历史
	b.history[userID] = append(b.history[userID], openai.ChatCompletionMessage{
		Role:    "user",
		Content: prompt,
	})

	// 3. 限制记忆长度（比如只记得最近 10 轮），防止 Token 爆炸
	if len(b.history[userID]) > 11 { // 1个系统提示词 + 10个对话
		b.history[userID] = append(b.history[userID][:1], b.history[userID][len(b.history[userID])-10:]...)
	}

	// 4. 发送完整的历史记录给大模型
	resp, _ := b.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    "deepseek-ai/DeepSeek-V3",
			Messages: b.history[userID],
		},
	)

	reply := resp.Choices[0].Message.Content

	// 5. 将 AI 的回答也存入历史
	b.history[userID] = append(b.history[userID], openai.ChatCompletionMessage{
		Role:    "assistant",
		Content: reply,
	})

	return reply
}