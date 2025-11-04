package services

import (
	"fmt"
	"strings"

	"github.com/Arterning/go-docx"
	"github.com/Arterning/go-rag/database"
	"github.com/Arterning/go-rag/models"
	"github.com/Arterning/go-rag/utils"
)

// ParseDocxFile 解析 docx 文件，返回标题和内容
func ParseDocxFile(filepath string) (title string, content string, err error) {
	// 使用 ExtractText 函数直接提取文本
	text, err := docx.ExtractText(filepath)
	if err != nil {
		return "", "", fmt.Errorf("打开 docx 文件失败: %w", err)
	}

	// 清理文本
	text = strings.TrimSpace(text)
	if text == "" {
		return "", "", fmt.Errorf("文档为空")
	}

	// 按行分割文本
	lines := strings.Split(text, "\n")

	// 第一个非空行作为标题
	title = ""
	contentStartIndex := 0
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			title = line
			contentStartIndex = i + 1
			break
		}
	}

	// 如果没有找到标题，使用默认标题
	if title == "" {
		title = "未命名文档"
		content = text
	} else {
		// 剩余部分作为内容
		if contentStartIndex < len(lines) {
			contentLines := lines[contentStartIndex:]
			content = strings.Join(contentLines, "\n")
			content = strings.TrimSpace(content)
		} else {
			// 如果只有标题，将标题也作为内容
			content = title
		}
	}

	return title, content, nil
}

// SaveDocument 保存文档及其分块到数据库
func SaveDocument(title, filename, content string) (*models.Document, error) {
	db := database.GetDB()

	// 创建文档记录
	doc := &models.Document{
		Title:    title,
		Filename: filename,
	}

	// 分块处理内容
	chunks := utils.ChunkText(content, utils.DefaultChunkSize, utils.DefaultOverlap)

	// 创建分块记录
	for i, chunkText := range chunks {
		chunk := models.DocumentChunk{
			ChunkText:  chunkText,
			ChunkIndex: i,
		}
		doc.Chunks = append(doc.Chunks, chunk)
	}

	// 保存到数据库
	if err := db.Create(doc).Error; err != nil {
		return nil, fmt.Errorf("保存文档失败: %w", err)
	}

	return doc, nil
}

// GetAllDocumentChunks 获取所有文档块（用于 RAG）
func GetAllDocumentChunks() ([]models.DocumentChunk, error) {
	db := database.GetDB()
	var chunks []models.DocumentChunk

	if err := db.Order("document_id, chunk_index").Find(&chunks).Error; err != nil {
		return nil, fmt.Errorf("获取文档块失败: %w", err)
	}

	return chunks, nil
}

// SearchChunks 简单的关键词搜索（用于 RAG 检索）
func SearchChunks(query string, limit int) ([]models.DocumentChunk, error) {
	db := database.GetDB()
	var chunks []models.DocumentChunk

	if limit <= 0 {
		limit = 5
	}

	// 简单的关键词匹配
	keywords := strings.Fields(query)
	if len(keywords) == 0 {
		return GetAllDocumentChunks()
	}

	// 使用 LIKE 查询
	queryStr := "%" + strings.Join(keywords, "%") + "%"

	if err := db.Where("chunk_text LIKE ?", queryStr).
		Order("document_id, chunk_index").
		Limit(limit).
		Find(&chunks).Error; err != nil {
		return nil, fmt.Errorf("搜索文档块失败: %w", err)
	}

	// 如果没有找到匹配的块，返回所有块（简单策略）
	if len(chunks) == 0 {
		return GetAllDocumentChunks()
	}

	return chunks, nil
}
