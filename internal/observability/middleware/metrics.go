package middleware

import (
	"fmt"
	"time"

	"task-manager/internal/observability/metrics"
	"task-manager/internal/observability/tracing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// MetricsMiddleware records HTTP metrics
func MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Generate and store trace ID
		traceID := uuid.New().String()
		tracer := tracing.NewSimpleTracer(traceID)
		c.Set("trace_id", traceID)
		c.Set("tracer", tracer)

		// Add trace ID to response header
		c.Header("X-Trace-ID", traceID)

		tracer.StartSpan(fmt.Sprintf("%s %s", c.Request.Method, c.Request.URL.Path))

		c.Next()

		duration := time.Since(start).Seconds()
		endpoint := c.Request.URL.Path

		// Record metrics
		metrics.RecordRequest(c.Request.Method, endpoint, c.Writer.Status(), duration)

		tracer.EndSpan()
	}
}

// GetTraceID retrieves the trace ID from gin context
func GetTraceID(c *gin.Context) string {
	if traceID, exists := c.Get("trace_id"); exists {
		return traceID.(string)
	}
	return ""
}

// GetTracer retrieves the tracer from gin context
func GetTracer(c *gin.Context) *tracing.SimpleTracer {
	if tracer, exists := c.Get("tracer"); exists {
		return tracer.(*tracing.SimpleTracer)
	}
	return nil
}
