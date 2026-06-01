package handler

import (
	"testing"

	"github.com/DilankaHer/sop-in-go/internal/app"
)

func TestNewHandlerWiresServerHandler(t *testing.T) {
	application := &app.App{}

	handler := NewHandler(application)

	if handler == nil {
		t.Fatal("handler is nil")
	}
	if handler.ServerHandler == nil {
		t.Fatal("ServerHandler is nil")
	}
}
