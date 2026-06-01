package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DilankaHer/sop-in-go/internal/app"
)

type versionResponseEnvelope struct {
	Status  int             `json:"status"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
	Error   json.RawMessage `json:"error"`
}

func TestVersionWritesCurrentVersion(t *testing.T) {
	application := &app.App{}
	handler := NewServerHandler(application)
	rec := httptest.NewRecorder()

	handler.Version(rec, httptest.NewRequest(http.MethodGet, "/version", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if got := rec.Header().Get("Content-Type"); got != "application/json" {
		t.Fatalf("Content-Type = %q, want %q", got, "application/json")
	}

	var got versionResponseEnvelope
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if got.Status != http.StatusOK {
		t.Fatalf("envelope status = %d, want %d", got.Status, http.StatusOK)
	}
	if got.Message != "success" {
		t.Fatalf("message = %q, want %q", got.Message, "success")
	}
	if string(got.Data) != `{"version":"v1.1"}` {
		t.Fatalf("data = %s, want %s", got.Data, `{"version":"v1.1"}`)
	}
	if string(got.Error) != "null" {
		t.Fatalf("error = %s, want null", got.Error)
	}
}

func TestNewServerHandlerStoresApp(t *testing.T) {
	application := &app.App{}

	handler := NewServerHandler(application)

	if handler == nil {
		t.Fatal("handler is nil")
	}
	if handler.app != application {
		t.Fatal("handler.app was not set from constructor argument")
	}
}
