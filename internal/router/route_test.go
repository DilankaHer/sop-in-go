package router

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DilankaHer/sop-in-go/internal/app"
	"github.com/DilankaHer/sop-in-go/internal/logger"
	"github.com/DilankaHer/sop-in-go/internal/middleware"
)

type routeResponseEnvelope struct {
	Status  int             `json:"status"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
	Error   json.RawMessage `json:"error"`
}

func TestInitRoutesRegistersVersion(t *testing.T) {
	log := logger.NewLogger()
	application := &app.App{
		Logger:     log,
		Middleware: middleware.NewMiddleware(log),
	}
	router := InitRoutes(application)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/version", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if got := rec.Header().Get("Content-Type"); got != "application/json" {
		t.Fatalf("Content-Type = %q, want %q", got, "application/json")
	}

	var got routeResponseEnvelope
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("failed to decode response body %q: %v", rec.Body.String(), err)
	}
	if got.Status != http.StatusOK {
		t.Fatalf("envelope status = %d, want %d", got.Status, http.StatusOK)
	}
	if got.Message != "OK" {
		t.Fatalf("message = %q, want %q", got.Message, "OK")
	}
	if string(got.Data) != "{\"version\":\"v1.1\"}" {
		t.Fatalf("data = %s, want %s", got.Data, "{\"version\":\"v1.1\"}")
	}
	if string(got.Error) != "null" {
		t.Fatalf("error = %s, want null", got.Error)
	}
}
