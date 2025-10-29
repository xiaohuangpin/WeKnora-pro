package middleware

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"

	"github.com/Tencent/WeKnora/internal/tracing"
	"github.com/Tencent/WeKnora/internal/types"
)

// Custom ResponseWriter to capture response content
type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// Override Write method to write response content to buffer and original writer
func (r responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

// TracingMiddleware provides a Gin middleware that creates a trace span for each request
func TracingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract trace context from request headers
		propagator := tracing.GetTracer()
		if propagator == nil {
			c.Next()
			return
		}

		// Get request ID as Span ID
		requestID := c.GetString(string(types.RequestIDContextKey))
		if requestID == "" {
			requestID = c.GetHeader("X-Request-ID")
		}

		// Create new span
		spanName := fmt.Sprintf("%s %s", c.Request.Method, c.FullPath())
		ctx, span := tracing.ContextWithSpan(c.Request.Context(), spanName)
		defer span.End()

		// Set basic span attributes
		span.SetAttributes(
			attribute.String("http.method", c.Request.Method),
			attribute.String("http.url", c.Request.URL.String()),
			attribute.String("http.path", c.FullPath()),
		)

		// Record request headers (optional, or selectively record important headers)
		for key, values := range c.Request.Header {
			// Skip sensitive or unnecessary headers
			if strings.ToLower(key) == "authorization" || strings.ToLower(key) == "cookie" {
				continue
			}
			span.SetAttributes(attribute.String("http.request.header."+key, strings.Join(values, ";")))
		}

		// Record request body (for POST/PUT/PATCH requests)
		if c.Request.Method == "POST" || c.Request.Method == "PUT" || c.Request.Method == "PATCH" {
			if c.Request.Body != nil {
				bodyBytes, _ := io.ReadAll(c.Request.Body)
				span.SetAttributes(attribute.String("http.request.body", string(bodyBytes)))
				// Reset request body because ReadAll consumes the Reader content
				c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			}
		}

		// Record query parameters
		if len(c.Request.URL.RawQuery) > 0 {
			span.SetAttributes(attribute.String("http.request.query", c.Request.URL.RawQuery))
		}

		// Set request context with span context
		c.Request = c.Request.WithContext(ctx)

		// Store tracing context in Gin context
		c.Set("trace.span", span)
		c.Set("trace.ctx", ctx)

		// Create response body capturer
		responseBody := &bytes.Buffer{}
		responseWriter := &responseBodyWriter{
			ResponseWriter: c.Writer,
			body:           responseBody,
		}
		c.Writer = responseWriter

		// Process request
		c.Next()

		// Set response status code
		statusCode := c.Writer.Status()
		span.SetAttributes(attribute.Int("http.status_code", statusCode))

		// Record response body
		responseContent := responseBody.String()
		if len(responseContent) > 0 {
			span.SetAttributes(attribute.String("http.response.body", responseContent))
		}

		// Record response headers (optional, or selectively record important headers)
		for key, values := range c.Writer.Header() {
			span.SetAttributes(attribute.String("http.response.header."+key, strings.Join(values, ";")))
		}

		// Mark as error if status code >= 400
		if statusCode >= 400 {
			span.SetStatus(codes.Error, fmt.Sprintf("HTTP %d", statusCode))
			if err := c.Errors.Last(); err != nil {
				span.RecordError(err.Err)
			}
		} else {
			span.SetStatus(codes.Ok, "")
		}
	}
}
