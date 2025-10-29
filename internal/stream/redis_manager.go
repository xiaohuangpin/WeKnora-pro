package stream

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Tencent/WeKnora/internal/types"
	"github.com/Tencent/WeKnora/internal/types/interfaces"
	"github.com/redis/go-redis/v9"
)

// redisStreamInfo Redis存储的流信息
type redisStreamInfo struct {
	SessionID           string           `json:"session_id"`
	RequestID           string           `json:"request_id"`
	Query               string           `json:"query"`
	Content             string           `json:"content"`
	KnowledgeReferences types.References `json:"knowledge_references"`
	LastUpdated         time.Time        `json:"last_updated"`
	IsCompleted         bool             `json:"is_completed"`
}

// RedisStreamManager 基于Redis的流管理器实现
type RedisStreamManager struct {
	client *redis.Client
	ttl    time.Duration // 流数据在Redis中的过期时间
	prefix string        // Redis键前缀
}

// NewRedisStreamManager 创建一个新的Redis流管理器
func NewRedisStreamManager(redisAddr, redisPassword string,
	redisDB int, prefix string, ttl time.Duration,
) (*RedisStreamManager, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       redisDB,
	})

	// 验证连接
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, fmt.Errorf("连接Redis失败: %w", err)
	}

	if ttl == 0 {
		ttl = 24 * time.Hour // 默认TTL为24小时
	}

	if prefix == "" {
		prefix = "stream:" // 默认前缀
	}

	return &RedisStreamManager{
		client: client,
		ttl:    ttl,
		prefix: prefix,
	}, nil
}

// 构建Redis键
func (r *RedisStreamManager) buildKey(sessionID, requestID string) string {
	return fmt.Sprintf("%s:%s:%s", r.prefix, sessionID, requestID)
}

// RegisterStream 注册一个新的流
func (r *RedisStreamManager) RegisterStream(ctx context.Context, sessionID, requestID, query string) error {
	info := &redisStreamInfo{
		SessionID:   sessionID,
		RequestID:   requestID,
		Query:       query,
		LastUpdated: time.Now(),
	}

	data, err := json.Marshal(info)
	if err != nil {
		return fmt.Errorf("序列化流信息失败: %w", err)
	}

	key := r.buildKey(sessionID, requestID)
	return r.client.Set(ctx, key, data, r.ttl).Err()
}

// UpdateStream 更新流内容
func (r *RedisStreamManager) UpdateStream(ctx context.Context, sessionID, requestID string, content string, references types.References) error {
	key := r.buildKey(sessionID, requestID)

	// 获取当前数据
	data, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil // 键不存在，可能已过期
		}
		return fmt.Errorf("获取流数据失败: %w", err)
	}

	var info redisStreamInfo
	if err := json.Unmarshal(data, &info); err != nil {
		return fmt.Errorf("解析流数据失败: %w", err)
	}

	// 更新数据
	info.Content += content
	if len(references) > 0 {
		info.KnowledgeReferences = references
	}
	info.LastUpdated = time.Now()

	// 保存回Redis
	updatedData, err := json.Marshal(info)
	if err != nil {
		return fmt.Errorf("序列化更新的流信息失败: %w", err)
	}

	return r.client.Set(ctx, key, updatedData, r.ttl).Err()
}

// CompleteStream 完成流
func (r *RedisStreamManager) CompleteStream(ctx context.Context, sessionID, requestID string) error {
	key := r.buildKey(sessionID, requestID)

	// 获取当前数据
	data, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil // 键不存在，可能已过期
		}
		return fmt.Errorf("获取流数据失败: %w", err)
	}

	var info redisStreamInfo
	if err := json.Unmarshal(data, &info); err != nil {
		return fmt.Errorf("解析流数据失败: %w", err)
	}

	// 标记为完成
	info.IsCompleted = true
	info.LastUpdated = time.Now()

	// 保存回Redis
	updatedData, err := json.Marshal(info)
	if err != nil {
		return fmt.Errorf("序列化更新的流信息失败: %w", err)
	}

	// 30s 后删除流
	go func() {
		time.Sleep(30 * time.Second)
		r.client.Del(ctx, key)
	}()
	return r.client.Set(ctx, key, updatedData, r.ttl).Err()
}

// GetStream 获取特定流
func (r *RedisStreamManager) GetStream(ctx context.Context, sessionID, requestID string) (*interfaces.StreamInfo, error) {
	key := r.buildKey(sessionID, requestID)

	// 获取数据
	data, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // 键不存在
		}
		return nil, fmt.Errorf("获取流数据失败: %w", err)
	}

	var info redisStreamInfo
	if err := json.Unmarshal(data, &info); err != nil {
		return nil, fmt.Errorf("解析流数据失败: %w", err)
	}

	// 转换为接口结构
	return &interfaces.StreamInfo{
		SessionID:           info.SessionID,
		RequestID:           info.RequestID,
		Query:               info.Query,
		Content:             info.Content,
		KnowledgeReferences: info.KnowledgeReferences,
		LastUpdated:         info.LastUpdated,
		IsCompleted:         info.IsCompleted,
	}, nil
}

// Close 关闭Redis连接
func (r *RedisStreamManager) Close() error {
	return r.client.Close()
}

// 确保实现了接口
var _ interfaces.StreamManager = (*RedisStreamManager)(nil)
