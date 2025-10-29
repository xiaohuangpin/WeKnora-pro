package stream

import (
	"context"
	"sync"
	"time"

	"github.com/Tencent/WeKnora/internal/types"
	"github.com/Tencent/WeKnora/internal/types/interfaces"
)

// 内存流信息
type memoryStreamInfo struct {
	sessionID           string
	requestID           string
	query               string
	content             string
	knowledgeReferences types.References
	lastUpdated         time.Time
	isCompleted         bool
}

// MemoryStreamManager 基于内存的流管理器实现
type MemoryStreamManager struct {
	// 会话ID -> 请求ID -> 流数据
	activeStreams map[string]map[string]*memoryStreamInfo
	mu            sync.RWMutex
}

// NewMemoryStreamManager 创建一个新的内存流管理器
func NewMemoryStreamManager() *MemoryStreamManager {
	return &MemoryStreamManager{
		activeStreams: make(map[string]map[string]*memoryStreamInfo),
	}
}

// RegisterStream 注册一个新的流
func (m *MemoryStreamManager) RegisterStream(ctx context.Context, sessionID, requestID, query string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	info := &memoryStreamInfo{
		sessionID:   sessionID,
		requestID:   requestID,
		query:       query,
		lastUpdated: time.Now(),
	}

	if _, exists := m.activeStreams[sessionID]; !exists {
		m.activeStreams[sessionID] = make(map[string]*memoryStreamInfo)
	}

	m.activeStreams[sessionID][requestID] = info
	return nil
}

// UpdateStream 更新流内容
func (m *MemoryStreamManager) UpdateStream(ctx context.Context,
	sessionID, requestID string, content string, references types.References,
) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if sessionMap, exists := m.activeStreams[sessionID]; exists {
		if stream, found := sessionMap[requestID]; found {
			stream.content += content
			if len(references) > 0 {
				stream.knowledgeReferences = references
			}
			stream.lastUpdated = time.Now()
		}
	}
	return nil
}

// CompleteStream 完成流
func (m *MemoryStreamManager) CompleteStream(ctx context.Context, sessionID, requestID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if sessionMap, exists := m.activeStreams[sessionID]; exists {
		if stream, found := sessionMap[requestID]; found {
			stream.isCompleted = true
			// 30s 后删除流
			go func() {
				time.Sleep(30 * time.Second)
				m.mu.Lock()
				defer m.mu.Unlock()
				delete(sessionMap, requestID)
				if len(sessionMap) == 0 {
					delete(m.activeStreams, sessionID)
				}
			}()
		}
	}
	return nil
}

// GetStream 获取特定流
func (m *MemoryStreamManager) GetStream(ctx context.Context,
	sessionID, requestID string,
) (*interfaces.StreamInfo, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if sessionMap, exists := m.activeStreams[sessionID]; exists {
		if stream, found := sessionMap[requestID]; found {
			return &interfaces.StreamInfo{
				SessionID:           stream.sessionID,
				RequestID:           stream.requestID,
				Query:               stream.query,
				Content:             stream.content,
				KnowledgeReferences: stream.knowledgeReferences,
				LastUpdated:         stream.lastUpdated,
				IsCompleted:         stream.isCompleted,
			}, nil
		}
	}
	return nil, nil
}

// 确保实现了接口
var _ interfaces.StreamManager = (*MemoryStreamManager)(nil)
