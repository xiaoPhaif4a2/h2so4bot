package llm

import (
	"bytes"
	"encoding/json"
	"h2so4bot/internal/database"
	"io"
	"log"
	"math/rand"
	"net/http"
	"time"
)

const (
	SiliconFlowURL = "https://api.siliconflow.cn/v1/chat/completions"
	ModelName      = "deepseek-ai/DeepSeek-V3"
)

type Brain struct {
	APIKey string
	Client *http.Client
}

func NewBrain(apiKey string) *Brain {
	return &Brain{
		APIKey: apiKey,
		// 🟢 将 30 * time.Second 改为 90 * time.Second
		Client: &http.Client{Timeout: 90 * time.Second},
	}
}

// 定义 API 请求结构体
type ChatRequest struct {
	Model            string        `json:"model"`
	Messages         []ChatMessage `json:"messages"`
	MaxTokens        int           `json:"max_tokens"`
	Temperature      float32       `json:"temperature"`       // 🟢 调高这个，让她更放飞自我
	FrequencyPenalty float32       `json:"frequency_penalty"` // 🟢 调高这个，阻止她反复说同样的话
}

// ChatMessage 就是之前报错找不到的那个结构体
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// 定义 API 响应结构体
type ChatResponse struct {
	Choices []struct {
		Message ChatMessage `json:"message"`
	} `json:"choices"`
}

// 🟢 注意这里加入了 groupID 参数
func (b *Brain) Think(sysPrompt string, userID, groupID int64, currentMsg string) string {
	// 1. 组装消息：先放入系统人设
	messages := []ChatMessage{
		{Role: "system", Content: sysPrompt},
	}

	// 2. 🟢 唤醒记忆：从数据库提取该群最近的 5 条聊天记录
	// (如果你之前的 db.go 里没写 GetContext，请告诉我，我发给你)
	history := database.GetContext(userID, groupID, 10)
	for _, h := range history {
		messages = append(messages, ChatMessage{Role: h.Role, Content: h.Content})
	}

	// 3. 压入当前最新消息
	messages = append(messages, ChatMessage{Role: "user", Content: currentMsg})

	// 4. 构造请求 (加入我们调优过的参数)
	reqBody := ChatRequest{
		Model:            ModelName,
		Messages:         messages,
		MaxTokens:        50,
		Temperature:      1.3,
		FrequencyPenalty: 0.5,
	}
	jsonData, _ := json.Marshal(reqBody)

	// ... 下面的发送 HTTP 请求和解析响应的代码保持完全不变 ...

	req, err := http.NewRequest("POST", SiliconFlowURL, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("⚠️ 构建请求失败: %v", err)
		return "哼，今天不想理你。"
	}

	req.Header.Set("Authorization", "Bearer "+b.APIKey)
	req.Header.Set("Content-Type", "application/json")

	// 3. 发送请求
	// 3. 发送请求
	resp, err := b.Client.Do(req)
	if err != nil {
		log.Printf("⚠️ 网络层断开: %v", err)
		return randomFallback()
	}
	defer resp.Body.Close() // 确保一定要关闭 Body

	// 🟢 如果状态码不是 200 OK，把硅基流动给的真实报错信息打印出来
	if resp.StatusCode != 200 {
		errorBody, _ := io.ReadAll(resp.Body)
		log.Printf("⚠️ API 拒绝了请求！状态码: %d, 详细原因: %s", resp.StatusCode, string(errorBody))
		return randomFallback()
	}

	// 4. 解析正常响应
	body, _ := io.ReadAll(resp.Body)
	var chatResp ChatResponse
	if err := json.Unmarshal(body, &chatResp); err != nil || len(chatResp.Choices) == 0 {
		log.Printf("⚠️ 解析响应失败: %v", err)
		return randomFallback()
	}

	return chatResp.Choices[0].Message.Content
}

func randomFallback() string {
	fallbacks := []string{
		"……网络卡了，不想理你。",
		"哼，凡人的服务器又崩溃了？",
		"（假装没听见）",
		"笨蛋，连句话都传不过来吗？",
		"本小姐现在累了，稍后再试吧。",
	}
	return fallbacks[rand.Intn(len(fallbacks))]
}
