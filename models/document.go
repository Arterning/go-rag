package models

import (
	"time"
)

// Document 代表上传的文档
type Document struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Title      string    `gorm:"type:text;not null" json:"title"`
	Filename   string    `gorm:"type:text;not null" json:"filename"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Chunks     []DocumentChunk `gorm:"foreignKey:DocumentID;constraint:OnDelete:CASCADE" json:"chunks,omitempty"`
}

// DocumentChunk 代表文档的一个分块
type DocumentChunk struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	DocumentID uint      `gorm:"not null;index" json:"document_id"`
	ChunkText  string    `gorm:"type:text;not null" json:"chunk_text"`
	ChunkIndex int       `gorm:"not null" json:"chunk_index"` // 块的顺序索引
	CreatedAt  time.Time `json:"created_at"`
}
