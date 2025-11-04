package services

import (
	"context"
	"fmt"
	"os"

	gollm "github.com/Arterning/go-llm"
)

var llmClient gollm.Client

// InitLLM 初始化 LLM 客户端
func InitLLM() error {
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		return fmt.Errorf("ANTHROPIC_API_KEY 环境变量未设置")
	}

	config := gollm.Config{
		Provider: gollm.ProviderClaude,
		Model:    "claude-3-5-sonnet-20241022", // 使用 Claude 3.5 Sonnet
		APIKey:   apiKey,
	}

	var err error
	llmClient, err = gollm.NewClient(config)
	if err != nil {
		return fmt.Errorf("初始化 LLM 客户端失败: %w", err)
	}

	return nil
}

// ChatWithContext 使用提供的上下文与 LLM 对话
func ChatWithContext(systemPrompt, userQuestion string) (string, error) {
	if llmClient == nil {
		return "", fmt.Errorf("LLM 客户端未初始化")
	}

	messages := []gollm.Message{
		gollm.NewTextMessage(gollm.RoleSystem, systemPrompt),
		gollm.NewTextMessage(gollm.RoleUser, userQuestion),
	}

	req := &gollm.ChatRequest{
		Messages:    messages,
		Temperature: 0.7,
		MaxTokens:   4096,
	}

	resp, err := llmClient.Chat(context.Background(), req)
	if err != nil {
		return "", fmt.Errorf("LLM 调用失败: %w", err)
	}

	if len(resp.Choices) > 0 && len(resp.Choices[0].Message.Content) > 0 {
		return resp.Choices[0].Message.Content[0].Text, nil
	}

	return "", fmt.Errorf("LLM 返回空响应")
}
