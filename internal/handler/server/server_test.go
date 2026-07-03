package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DilankaHer/sop-in-go/internal/app"
)

func TestVersionReturnsCurrentVersion(t *testing.T) {
	handler := NewServerHandler(&app.App{})

	resp := handler.Version(httptest.NewRequest(http.MethodGet, "/version", nil))

	if resp.Status != http.StatusOK {
		t.Fatalf("status = %d, want %d", resp.Status, http.StatusOK)
	}
	if resp.Message != http.StatusText(http.StatusOK) {
		t.Fatalf("message = %q, want %q", resp.Message, http.StatusText(http.StatusOK))
	}
	data, ok := resp.Data.(getVersionResp)
	if !ok {
		t.Fatalf("data type = %T, want getVersionResp", resp.Data)
	}
	if data.Version != "v1.1" {
		t.Fatalf("version = %q, want %q", data.Version, "v1.1")
	}
	if resp.Error != nil {
		t.Fatalf("error = %v, want nil", resp.Error)
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
