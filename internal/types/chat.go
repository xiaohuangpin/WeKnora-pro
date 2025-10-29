package types

import (
	"database/sql/driver"
	"encoding/json"
)

// ChatResponse chat response
type ChatResponse struct {
	Content string `json:"content"`
	// Usage information
	Usage struct {
		// Prompt tokens
		PromptTokens int `json:"prompt_tokens"`
		// Completion tokens
		CompletionTokens int `json:"completion_tokens"`
		// Total tokens
		TotalTokens int `json:"total_tokens"`
	} `json:"usage"`
}

// Response type
type ResponseType string

const (
	// Answer response type
	ResponseTypeAnswer ResponseType = "answer"
	// References response type
	ResponseTypeReferences ResponseType = "references"
)

// StreamResponse stream response
type StreamResponse struct {
	// Unique identifier
	ID string `json:"id"`
	// Response type
	ResponseType ResponseType `json:"response_type"`
	// Current fragment content
	Content string `json:"content"`
	// Whether the response is complete
	Done bool `json:"done"`
	// Knowledge references
	KnowledgeReferences References `json:"knowledge_references"`
}

// References references
type References []*SearchResult

// Value implements the driver.Valuer interface, used to convert References to database values
func (c References) Value() (driver.Value, error) {
	return json.Marshal(c)
}

// Scan implements the sql.Scanner interface, used to convert database values to References
func (c *References) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	b, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(b, c)
}
