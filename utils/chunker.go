package utils

import (
	"strings"
	"unicode/utf8"
)

const (
	// DefaultChunkSize 默认每块大小（字符数）
	DefaultChunkSize = 1000
	// DefaultOverlap 默认重叠大小（字符数）
	DefaultOverlap = 200
)

// ChunkText 将文本分块，支持重叠
func ChunkText(text string, chunkSize, overlap int) []string {
	if chunkSize <= 0 {
		chunkSize = DefaultChunkSize
	}
	if overlap < 0 {
		overlap = DefaultOverlap
	}
	if overlap >= chunkSize {
		overlap = chunkSize / 4 // 防止重叠过大
	}

	// 清理文本
	text = strings.TrimSpace(text)
	if text == "" {
		return []string{}
	}

	textLen := utf8.RuneCountInString(text)
	if textLen <= chunkSize {
		return []string{text}
	}

	var chunks []string
	runes := []rune(text)
	start := 0

	for start < textLen {
		end := start + chunkSize
		if end > textLen {
			end = textLen
		}

		// 尝试在句子边界处分块
		chunk := string(runes[start:end])

		// 如果不是最后一块，尝试在句号、问号、感叹号等处断开
		if end < textLen {
			if lastPeriod := strings.LastIndexAny(chunk, "。！？\n.!?"); lastPeriod != -1 && lastPeriod > chunkSize/2 {
				chunk = chunk[:lastPeriod+1]
				end = start + utf8.RuneCountInString(chunk)
			}
		}

		chunks = append(chunks, strings.TrimSpace(chunk))

		// 下一块的开始位置（考虑重叠）
		start = end - overlap
		if start < 0 {
			start = 0
		}
	}

	return chunks
}
