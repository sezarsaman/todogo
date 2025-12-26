package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

// MockTaskHandler for testing
type MockTaskHandler struct{}

func (m *MockTaskHandler) Register(r *gin.Engine) {}

func TestNewRouter(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("health_endpoint_returns_200", func(t *testing.T) {
		handler := &MockTaskHandler{}
		router := NewRouter(handler)

		req := httptest.NewRequest(http.MethodGet, "/health", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
	})

	t.Run("health_endpoint_returns_ok_message", func(t *testing.T) {
		handler := &MockTaskHandler{}
		router := NewRouter(handler)

		req := httptest.NewRequest(http.MethodGet, "/health", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Body.String() != "{\"message\":\"OK\"}" {
			t.Errorf("unexpected response body: %s", w.Body.String())
		}
	})

	t.Run("undefined_route_returns_404", func(t *testing.T) {
		handler := &MockTaskHandler{}
		router := NewRouter(handler)

		req := httptest.NewRequest(http.MethodGet, "/notfound", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("expected status 404, got %d", w.Code)
		}
	})
}
