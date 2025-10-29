package middleware

import (
	"context"
	"errors"
	"log"
	"net/http"
	"slices"
	"strings"

	"github.com/Tencent/WeKnora/internal/config"
	"github.com/Tencent/WeKnora/internal/types"
	"github.com/Tencent/WeKnora/internal/types/interfaces"
	"github.com/gin-gonic/gin"
)

// 无需认证的API列表
var noAuthAPI = map[string][]string{
	"/health":               {"GET"},
	"/api/v1/auth/register": {"POST"},
	"/api/v1/auth/login":    {"POST"},
	"/api/v1/auth/refresh":  {"POST"},
}

// 检查请求是否在无需认证的API列表中
func isNoAuthAPI(path string, method string) bool {
	for api, methods := range noAuthAPI {
		// 如果以*结尾，按照前缀匹配，否则按照全路径匹配
		if strings.HasSuffix(api, "*") {
			if strings.HasPrefix(path, strings.TrimSuffix(api, "*")) && slices.Contains(methods, method) {
				return true
			}
		} else if path == api && slices.Contains(methods, method) {
			return true
		}
	}
	return false
}

// Auth 认证中间件
func Auth(tenantService interfaces.TenantService, userService interfaces.UserService, cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// ignore OPTIONS request
		if c.Request.Method == "OPTIONS" {
			c.Next()
			return
		}

		// 检查请求是否在无需认证的API列表中
		if isNoAuthAPI(c.Request.URL.Path, c.Request.Method) {
			c.Next()
			return
		}

		// 尝试JWT Token认证
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			token := strings.TrimPrefix(authHeader, "Bearer ")
			user, err := userService.ValidateToken(c.Request.Context(), token)
			if err == nil && user != nil {
				// JWT Token认证成功
				// 获取租户信息
				tenant, err := tenantService.GetTenantByID(c.Request.Context(), user.TenantID)
				if err != nil {
					log.Printf("Error getting tenant by ID: %v, tenantID: %d, userID: %s", err, user.TenantID, user.ID)
					c.JSON(http.StatusUnauthorized, gin.H{
						"error": "Unauthorized: invalid tenant",
					})
					c.Abort()
					return
				}

				// 存储用户和租户信息到上下文
				c.Set(types.TenantIDContextKey.String(), user.TenantID)
				c.Set(types.TenantInfoContextKey.String(), tenant)
				c.Set("user", user)
				c.Request = c.Request.WithContext(
					context.WithValue(
						context.WithValue(
							context.WithValue(c.Request.Context(), types.TenantIDContextKey, user.TenantID),
							types.TenantInfoContextKey, tenant,
						),
						"user", user,
					),
				)
				c.Next()
				return
			}
		}

		// 尝试X-API-Key认证（兼容模式）
		apiKey := c.GetHeader("X-API-Key")
		if apiKey != "" {
			// Get tenant information
			tenantID, err := tenantService.ExtractTenantIDFromAPIKey(apiKey)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": "Unauthorized: invalid API key format",
				})
				c.Abort()
				return
			}

			// Verify API key validity (matches the one in database)
			t, err := tenantService.GetTenantByID(c.Request.Context(), tenantID)
			if err != nil {
				log.Printf("Error getting tenant by ID: %v, tenantID: %d, apiKey: %s", err, tenantID, apiKey)
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": "Unauthorized: invalid API key",
				})
				c.Abort()
				return
			}

			if t == nil || t.APIKey != apiKey {
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": "Unauthorized: invalid API key",
				})
				c.Abort()
				return
			}

			// Store tenant ID in context
			c.Set(types.TenantIDContextKey.String(), tenantID)
			c.Set(types.TenantInfoContextKey.String(), t)
			c.Request = c.Request.WithContext(
				context.WithValue(
					context.WithValue(c.Request.Context(), types.TenantIDContextKey, tenantID),
					types.TenantInfoContextKey, t,
				),
			)
			c.Next()
			return
		}

		// 没有提供任何认证信息
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: missing authentication"})
		c.Abort()
	}
}

// GetTenantIDFromContext helper function to get tenant ID from context
func GetTenantIDFromContext(ctx context.Context) (uint, error) {
	tenantID, ok := ctx.Value("tenantID").(uint)
	if !ok {
		return 0, errors.New("tenant ID not found in context")
	}
	return tenantID, nil
}
