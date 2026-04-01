package main

import (
	"context"
	"fmt"
	"github.com/sashabaranov/go-openai"
)

func main() {
	// 1. 配置硅基流动的 API 密钥和基础网址
	config := openai.DefaultConfig("sk-nkywmeisoxvorpxtqefbsxsnacvboodbecpjfhacjoxrehtp")
	config.BaseURL = "https://api.siliconflow.cn/v1" // 硅基流动的接口地址

	// 2. 创建一个客户端实例
	client := openai.NewClientWithConfig(config)

	// 3. 设定对话内容：这就像在写微小说的角色设定，告诉它扮演谁
	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: "你现在是一个傲娇的文学少女，说话总是喜欢引用一两句经典名著，但态度又有些口是心非。回答要尽量精简。", // 这里就是塑造人设的地方
		},
		{
			Role:    openai.ChatMessageRoleUser,
			Content: "初次见面，请多关照！", // 这是你对它说的话
		},
	}

	// 4. 发送请求给大模型
	fmt.Println("大脑正在思考中...")
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    "deepseek-ai/DeepSeek-V3", // 这里填硅基流动支持的模型名称，比如 Qwen/Qwen2.5-7B-Instruct
			Messages: messages,
		},
	)

	// 5. 处理结果
	if err != nil {
		fmt.Printf("调用 API 失败了: %v\n", err)
		return
	}

	fmt.Println("机器人回复:", resp.Choices[0].Message.Content)
}