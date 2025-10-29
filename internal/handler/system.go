package handler

import (
	"github.com/Tencent/WeKnora/internal/logger"
	"github.com/gin-gonic/gin"
)

// SystemHandler handles system-related requests
type SystemHandler struct{}

// NewSystemHandler creates a new system handler
func NewSystemHandler() *SystemHandler {
	return &SystemHandler{}
}

// GetSystemInfoResponse defines the response structure for system info
type GetSystemInfoResponse struct {
	Version   string `json:"version"`
	CommitID  string `json:"commit_id,omitempty"`
	BuildTime string `json:"build_time,omitempty"`
	GoVersion string `json:"go_version,omitempty"`
}

// 编译时注入的版本信息
var (
	Version   = "unknown"
	CommitID  = "unknown"
	BuildTime = "unknown"
	GoVersion = "unknown"
)

// GetSystemInfo gets system information including version and commit ID
func (h *SystemHandler) GetSystemInfo(c *gin.Context) {
	ctx := logger.CloneContext(c.Request.Context())

	response := GetSystemInfoResponse{
		Version:   Version,
		CommitID:  CommitID,
		BuildTime: BuildTime,
		GoVersion: GoVersion,
	}

	logger.Info(ctx, "System info retrieved successfully")
	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "success",
		"data": response,
	})
}
