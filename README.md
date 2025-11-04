# Go RAG 应用

一个基于 Go 的 RAG (Retrieval-Augmented Generation) 应用，支持上传 docx 文件并使用 AI 进行问答。

## 技术栈

- **Web 框架**: Gin
- **数据库**: SQLite (GORM)
- **文档解析**: go-docx
- **LLM 集成**: go-llm (支持多个 LLM 提供商)
- **AI 模型**: Claude 3.5 Sonnet

## 功能特性

1. **文档上传**: 上传 docx 文件，自动提取标题和内容
2. **文档分块**: 自动将文档内容分块存储，优化 RAG 检索
3. **AI 问答**: 基于上传的文档内容进行智能问答
4. **多 LLM 支持**: 使用 go-llm 库，可轻松切换不同的 LLM 提供商

## 快速开始

### 1. 环境准备

确保已安装 Go 1.23 或更高版本。

### 2. 配置环境变量

复制 `.env.example` 文件并重命名为 `.env`，然后配置您的 API Key：

```bash
cp .env.example .env
```

编辑 `.env` 文件，设置您的 Anthropic API Key：

```env
ANTHROPIC_API_KEY=your_actual_api_key_here
```

或者直接在命令行中设置环境变量：

**Windows (CMD)**:
```cmd
set ANTHROPIC_API_KEY=your_api_key_here
```

**Windows (PowerShell)**:
```powershell
$env:ANTHROPIC_API_KEY="your_api_key_here"
```

**Linux/Mac**:
```bash
export ANTHROPIC_API_KEY=your_api_key_here
```

### 3. 安装依赖

```bash
go mod download
```

### 4. 运行应用

```bash
go run main.go
```

服务器将在 `http://localhost:8080` 启动。

## API 接口

### 1. 健康检查

**请求**:
```
GET /api/health
```

**响应**:
```json
{
  "status": "ok",
  "message": "RAG 服务运行正常"
}
```

### 2. 上传文档

**请求**:
```
POST /api/upload
Content-Type: multipart/form-data

file: your_document.docx
```

**响应示例**:
```json
{
  "success": true,
  "message": "文件上传成功",
  "document_id": 1,
  "title": "文档标题",
  "chunk_count": 5
}
```

**使用 curl 测试**:
```bash
curl -X POST http://localhost:8080/api/upload \
  -F "file=@your_document.docx"
```

### 3. AI 问答

**请求**:
```
POST /api/qa
Content-Type: application/json

{
  "question": "你的问题"
}
```

**响应示例**:
```json
{
  "success": true,
  "answer": "根据文档内容，答案是..."
}
```

**使用 curl 测试**:
```bash
curl -X POST http://localhost:8080/api/qa \
  -H "Content-Type: application/json" \
  -d '{"question":"文档的主要内容是什么？"}'
```

## 项目结构

```
go-rag/
├── main.go              # 主程序入口
├── go.mod               # Go 模块配置
├── go.sum               # 依赖校验
├── .env.example         # 环境变量示例
├── README.md            # 项目说明
├── documents.db         # SQLite 数据库（自动生成）
├── temp/                # 临时文件目录（自动生成）
├── models/              # 数据模型
│   └── document.go      # 文档和文档块模型
├── database/            # 数据库配置
│   └── db.go            # 数据库初始化
├── handlers/            # HTTP 处理器
│   ├── upload.go        # 文件上传处理
│   └── qa.go            # 问答处理
├── services/            # 业务逻辑
│   ├── document_service.go  # 文档服务
│   ├── llm_service.go       # LLM 服务
│   └── rag_service.go       # RAG 问答服务
└── utils/               # 工具函数
    └── chunker.go       # 文本分块工具
```

## 配置说明

### 文档分块参数

在 `utils/chunker.go` 中可以调整分块参数：

- `DefaultChunkSize`: 每块的字符数（默认 1000）
- `DefaultOverlap`: 块之间的重叠字符数（默认 200）

### LLM 配置

在 `services/llm_service.go` 中可以修改：

- `Model`: 使用的模型（默认 `claude-3-5-sonnet-20241022`）
- `Temperature`: 生成温度（默认 0.7）
- `MaxTokens`: 最大 token 数（默认 4096）

## RAG 实现说明

本项目采用简单但有效的 RAG 策略：

1. **文档存储**: 将上传的 docx 文件分块存储到 SQLite 数据库
2. **检索策略**: 将所有文档块提供给 LLM（适合文档数量不多的场景）
3. **生成回答**: 使用 Claude 根据文档内容生成准确的回答

## 扩展功能建议

如果需要处理大量文档，可以考虑以下优化：

1. **向量搜索**: 集成向量数据库（如 chromem-go）实现语义搜索
2. **BM25 检索**: 使用 BM25 算法进行关键词匹配
3. **混合检索**: 结合关键词和语义搜索
4. **缓存机制**: 缓存常见问题的答案

## 常见问题

### Q: 如何获取 Anthropic API Key?
A: 访问 [Anthropic Console](https://console.anthropic.com/) 注册并创建 API Key。

### Q: 支持其他 LLM 提供商吗?
A: 支持！go-llm 库支持多个提供商（OpenAI、Claude、Gemini、DeepSeek 等）。修改 `services/llm_service.go` 中的配置即可。

### Q: 支持哪些文档格式?
A: 目前仅支持 .docx 格式。如需支持其他格式，可以添加相应的解析器。

### Q: 如何处理大型文档?
A: 系统会自动将文档分块。可以在 `utils/chunker.go` 中调整分块大小。

## 许可证

MIT License
