package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Arterning/go-rag/services"
	"github.com/gin-gonic/gin"
)

// UploadRequest 上传请求的响应
type UploadResponse struct {
	Success   bool   `json:"success"`
	Message   string `json:"message"`
	DocumentID uint  `json:"document_id,omitempty"`
	Title     string `json:"title,omitempty"`
	ChunkCount int  `json:"chunk_count,omitempty"`
}

// UploadDocxHandler 处理 docx 文件上传
func UploadDocxHandler(c *gin.Context) {
	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, UploadResponse{
			Success: false,
			Message: "未找到上传文件",
		})
		return
	}

	// 验证文件扩展名
	ext := filepath.Ext(file.Filename)
	if ext != ".docx" {
		c.JSON(http.StatusBadRequest, UploadResponse{
			Success: false,
			Message: "只支持 .docx 格式的文件",
		})
		return
	}

	// 创建临时目录
	tempDir := "./temp"
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, UploadResponse{
			Success: false,
			Message: "创建临时目录失败",
		})
		return
	}

	// 保存上传的文件到临时目录
	tempFilePath := filepath.Join(tempDir, file.Filename)
	if err := c.SaveUploadedFile(file, tempFilePath); err != nil {
		c.JSON(http.StatusInternalServerError, UploadResponse{
			Success: false,
			Message: "保存文件失败",
		})
		return
	}

	// 延迟删除临时文件
	defer os.Remove(tempFilePath)

	// 解析 docx 文件
	title, content, err := services.ParseDocxFile(tempFilePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, UploadResponse{
			Success: false,
			Message: fmt.Sprintf("解析文件失败: %v", err),
		})
		return
	}

	// 保存到数据库
	doc, err := services.SaveDocument(title, file.Filename, content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, UploadResponse{
			Success: false,
			Message: fmt.Sprintf("保存文档失败: %v", err),
		})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, UploadResponse{
		Success:    true,
		Message:    "文件上传成功",
		DocumentID: doc.ID,
		Title:      doc.Title,
		ChunkCount: len(doc.Chunks),
	})
}
