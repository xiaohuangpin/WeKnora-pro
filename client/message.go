// Package client provides the implementation for interacting with the WeKnora API
// The Message related interfaces are used to manage messages in a session
// Messages can be created, retrieved, deleted, and queried
package client

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// Message message information
type Message struct {
	ID                  string          `json:"id"`
	SessionID           string          `json:"session_id"`
	RequestID           string          `json:"request_id"`
	Content             string          `json:"content"`
	Role                string          `json:"role"`
	KnowledgeReferences []*SearchResult `json:"knowledge_references" `
	IsCompleted         bool            `json:"is_completed"`
	CreatedAt           time.Time       `json:"created_at"`
	UpdatedAt           time.Time       `json:"updated_at"`
}

// MessageListResponse message list response
type MessageListResponse struct {
	Success bool      `json:"success"`
	Data    []Message `json:"data"`
}

// LoadMessages loads session messages, supports pagination and time filtering
func (c *Client) LoadMessages(ctx context.Context, sessionID string, limit int, beforeTime *time.Time) ([]Message, error) {
	path := fmt.Sprintf("/api/v1/messages/%s/load", sessionID)

	queryParams := url.Values{}
	queryParams.Add("limit", strconv.Itoa(limit))

	if beforeTime != nil {
		queryParams.Add("before_time", beforeTime.Format(time.RFC3339Nano))
	}

	resp, err := c.doRequest(ctx, http.MethodGet, path, nil, queryParams)
	if err != nil {
		return nil, err
	}

	var response MessageListResponse
	if err := parseResponse(resp, &response); err != nil {
		return nil, err
	}

	return response.Data, nil
}

// GetRecentMessages gets recent messages from a session
func (c *Client) GetRecentMessages(ctx context.Context, sessionID string, limit int) ([]Message, error) {
	return c.LoadMessages(ctx, sessionID, limit, nil)
}

// GetMessagesBefore gets messages before a specified time
func (c *Client) GetMessagesBefore(ctx context.Context, sessionID string, beforeTime time.Time, limit int) ([]Message, error) {
	return c.LoadMessages(ctx, sessionID, limit, &beforeTime)
}

// DeleteMessage deletes a message
func (c *Client) DeleteMessage(ctx context.Context, sessionID string, messageID string) error {
	path := fmt.Sprintf("/api/v1/messages/%s/%s", sessionID, messageID)
	resp, err := c.doRequest(ctx, http.MethodDelete, path, nil, nil)
	if err != nil {
		return err
	}

	var response struct {
		Success bool   `json:"success"`
		Message string `json:"message,omitempty"`
	}

	return parseResponse(resp, &response)
}
