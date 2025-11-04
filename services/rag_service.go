package services

import (
	"fmt"
	"strings"
)

// RAGQuery 执行 RAG 问答
func RAGQuery(question string) (string, error) {
	// 1. 从数据库检索相关文档块
	chunks, err := GetAllDocumentChunks()
	if err != nil {
		return "", fmt.Errorf("检索文档失败: %w", err)
	}

	if len(chunks) == 0 {
		return "", fmt.Errorf("数据库中没有文档，请先上传文档")
	}

	// 2. 构建上下文
	// 简单策略：将所有文档块都提供给 AI
	var contextBuilder strings.Builder
	contextBuilder.WriteString("以下是所有可用的文档内容：\n\n")

	for i, chunk := range chunks {
		contextBuilder.WriteString(fmt.Sprintf("【文档块 %d】\n%s\n\n", i+1, chunk.ChunkText))
	}

	context := contextBuilder.String()

	// 3. 构建系统提示词
	systemPrompt := fmt.Sprintf(`你是一个专业的文档问答助手。你的任务是根据提供的文档内容回答用户的问题。

请遵循以下规则：
1. 仅根据提供的文档内容回答问题
2. 如果文档中没有相关信息，请明确告知用户
3. 回答要准确、简洁、有条理
4. 如果可能，引用具体的文档内容来支持你的回答

%s`, context)

	// 4. 调用 LLM 生成回答
	answer, err := ChatWithContext(systemPrompt, question)
	if err != nil {
		return "", fmt.Errorf("生成回答失败: %w", err)
	}

	return answer, nil
}

// RAGQueryWithSearch 使用关键词搜索的 RAG 问答（可选的优化版本）
func RAGQueryWithSearch(question string, maxChunks int) (string, error) {
	// 1. 使用关键词搜索相关文档块
	chunks, err := SearchChunks(question, maxChunks)
	if err != nil {
		return "", fmt.Errorf("搜索文档失败: %w", err)
	}

	if len(chunks) == 0 {
		return "", fmt.Errorf("数据库中没有文档，请先上传文档")
	}

	// 2. 构建上下文
	var contextBuilder strings.Builder
	contextBuilder.WriteString("以下是与您的问题相关的文档内容：\n\n")

	for i, chunk := range chunks {
		contextBuilder.WriteString(fmt.Sprintf("【文档块 %d】\n%s\n\n", i+1, chunk.ChunkText))
	}

	context := contextBuilder.String()

	// 3. 构建系统提示词
	systemPrompt := fmt.Sprintf(`你是一个专业的文档问答助手。你的任务是根据提供的文档内容回答用户的问题。

请遵循以下规则：
1. 仅根据提供的文档内容回答问题
2. 如果文档中没有相关信息，请明确告知用户
3. 回答要准确、简洁、有条理
4. 如果可能，引用具体的文档内容来支持你的回答

%s`, context)

	// 4. 调用 LLM 生成回答
	answer, err := ChatWithContext(systemPrompt, question)
	if err != nil {
		return "", fmt.Errorf("生成回答失败: %w", err)
	}

	return answer, nil
}
