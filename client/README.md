# WeKnora HTTP 客户端

这个包提供了与WeKnora服务进行交互的客户端库，支持所有基于HTTP的接口调用，使其他模块更方便地集成WeKnora服务，无需直接编写HTTP请求代码。

## 主要功能

该客户端包含以下主要功能模块：

1. **会话管理**：创建、获取、更新和删除会话
2. **知识库管理**：创建、获取、更新和删除知识库
3. **知识管理**：添加、获取和删除知识内容
4. **租户管理**：租户的CRUD操作
5. **知识问答**：支持普通问答和流式问答
6. **分块管理**：查询、更新和删除知识分块
7. **消息管理**：获取和删除会话消息
8. **模型管理**：创建、获取、更新和删除模型

## 使用方法

### 创建客户端实例

```go
import (
    "context"
    "github.com/Tencent/WeKnora/internal/client"
    "time"
)

// 创建客户端实例
apiClient := client.NewClient(
    "http://api.example.com", 
    client.WithToken("your-auth-token"),
    client.WithTimeout(30*time.Second),
)
```

### 示例：创建知识库并上传文件

```go
// 创建知识库
kb := &client.KnowledgeBase{
    Name:        "测试知识库",
    Description: "这是一个测试知识库",
    ChunkingConfig: client.ChunkingConfig{
        ChunkSize:    500,
        ChunkOverlap: 50,
        Separators:   []string{"\n\n", "\n", ". ", "? ", "! "},
    },
    ImageProcessingConfig: client.ImageProcessingConfig{
        ModelID: "image_model_id",
    },
    EmbeddingModelID: "embedding_model_id",
    SummaryModelID:   "summary_model_id",
}

kb, err := apiClient.CreateKnowledgeBase(context.Background(), kb)
if err != nil {
    // 处理错误
}

// 上传知识文件并添加元数据
metadata := map[string]string{
    "source": "local",
    "type":   "document",
}
knowledge, err := apiClient.CreateKnowledgeFromFile(context.Background(), kb.ID, "path/to/file.pdf", metadata)
if err != nil {
    // 处理错误
}
```

### 示例：创建会话并进行问答

```go
// 创建会话
sessionRequest := &client.CreateSessionRequest{
    KnowledgeBaseID: knowledgeBaseID,
    SessionStrategy: &client.SessionStrategy{
        MaxRounds:        10,
        EnableRewrite:    true,
        FallbackStrategy: "fixed_answer",
        FallbackResponse: "抱歉，我无法回答这个问题",
        EmbeddingTopK:    5,
        KeywordThreshold: 0.5,
        VectorThreshold:  0.7,
        RerankModelID:    "rerank_model_id",
        RerankTopK:       3,
        RerankThreshold:  0.8,
        SummaryModelID:   "summary_model_id",
    },
}

session, err := apiClient.CreateSession(context.Background(), sessionRequest)
if err != nil {
    // 处理错误
}

// 普通问答
answer, err := apiClient.KnowledgeQA(context.Background(), session.ID, &client.KnowledgeQARequest{
    Query: "什么是人工智能?",
})
if err != nil {
    // 处理错误
}

// 流式问答
err = apiClient.KnowledgeQAStream(context.Background(), session.ID, "什么是机器学习?", func(response *client.StreamResponse) error {
    // 处理每个响应片段
    fmt.Print(response.Content)
    return nil
})
if err != nil {
    // 处理错误
}
```

### 示例：管理模型

```go
// 创建模型
modelRequest := &client.CreateModelRequest{
    Name:        "测试模型",
    Type:        client.ModelTypeChat,
    Source:      client.ModelSourceInternal,
    Description: "这是一个测试模型",
    Parameters: client.ModelParameters{
        "temperature": 0.7,
        "top_p":       0.9,
    },
    IsDefault: true,
}
model, err := apiClient.CreateModel(context.Background(), modelRequest)
if err != nil {
    // 处理错误
}

// 列出所有模型
models, err := apiClient.ListModels(context.Background())
if err != nil {
    // 处理错误
}
```

### 示例：管理知识分块

```go
// 列出知识分块
chunks, total, err := apiClient.ListKnowledgeChunks(context.Background(), knowledgeID, 1, 10)
if err != nil {
    // 处理错误
}

// 更新分块
updateRequest := &client.UpdateChunkRequest{
    Content:   "更新后的分块内容",
    IsEnabled: true,
}
updatedChunk, err := apiClient.UpdateChunk(context.Background(), knowledgeID, chunkID, updateRequest)
if err != nil {
    // 处理错误
}
```

### 示例：获取会话消息

```go
// 获取最近消息
messages, err := apiClient.GetRecentMessages(context.Background(), sessionID, 10)
if err != nil {
    // 处理错误
}

// 获取指定时间之前的消息
beforeTime := time.Now().Add(-24 * time.Hour)
olderMessages, err := apiClient.GetMessagesBefore(context.Background(), sessionID, beforeTime, 10)
if err != nil {
    // 处理错误
}
```

## 完整示例

请参考 `example.go` 文件中的 `ExampleUsage` 函数，其中展示了客户端的完整使用流程。