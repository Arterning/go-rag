package handlers

import (
	"net/http"

	"github.com/Arterning/go-rag/services"
	"github.com/gin-gonic/gin"
)

// QARequest 问答请求
type QARequest struct {
	Question string `json:"question" binding:"required"`
}

// QAResponse 问答响应
type QAResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Answer  string `json:"answer,omitempty"`
}

// QAHandler 处理 AI 问答请求
func QAHandler(c *gin.Context) {
	var req QARequest

	// 解析请求体
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, QAResponse{
			Success: false,
			Message: "请求格式错误：question 字段为必填项",
		})
		return
	}

	// 验证问题不能为空
	if req.Question == "" {
		c.JSON(http.StatusBadRequest, QAResponse{
			Success: false,
			Message: "问题不能为空",
		})
		return
	}

	// 执行 RAG 问答
	answer, err := services.RAGQuery(req.Question)
	if err != nil {
		c.JSON(http.StatusInternalServerError, QAResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	// 返回回答
	c.JSON(http.StatusOK, QAResponse{
		Success: true,
		Answer:  answer,
	})
}
