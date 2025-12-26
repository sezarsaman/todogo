package tracing

import (
	"context"
	"fmt"
	"log"
)

// SimpleTracer provides basic request tracing
type SimpleTracer struct {
	traceID   string
	spanStack []string
}

// NewSimpleTracer creates a new tracer
func NewSimpleTracer(traceID string) *SimpleTracer {
	return &SimpleTracer{
		traceID:   traceID,
		spanStack: []string{},
	}
}

// StartSpan starts a new span
func (t *SimpleTracer) StartSpan(name string) {
	t.spanStack = append(t.spanStack, name)
	log.Printf("[TRACE] traceID=%s spanID=%d span.start name=%s", t.traceID, len(t.spanStack), name)
}

// EndSpan ends the current span
func (t *SimpleTracer) EndSpan() {
	if len(t.spanStack) > 0 {
		spanName := t.spanStack[len(t.spanStack)-1]
		t.spanStack = t.spanStack[:len(t.spanStack)-1]
		log.Printf("[TRACE] traceID=%s span.end name=%s", t.traceID, spanName)
	}
}

// RecordEvent records an event in the trace
func (t *SimpleTracer) RecordEvent(name string, attributes map[string]interface{}) {
	attrs := ""
	for k, v := range attributes {
		attrs += fmt.Sprintf(" %s=%v", k, v)
	}
	log.Printf("[TRACE] traceID=%s event=%s%s", t.traceID, name, attrs)
}

// GetTraceID returns the trace ID
func (t *SimpleTracer) GetTraceID() string {
	return t.traceID
}

// ContextKey for storing tracer in context
type ContextKey string

const TracerKey ContextKey = "tracer"

// WithTracer adds a tracer to the context
func WithTracer(ctx context.Context, tracer *SimpleTracer) context.Context {
	return context.WithValue(ctx, TracerKey, tracer)
}

// FromContext retrieves a tracer from context
func FromContext(ctx context.Context) *SimpleTracer {
	tracer, ok := ctx.Value(TracerKey).(*SimpleTracer)
	if !ok {
		return nil
	}
	return tracer
}
