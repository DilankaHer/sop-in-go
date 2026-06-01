package utility

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

type responseEnvelope struct {
	Status  int             `json:"status"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
	Error   json.RawMessage `json:"error"`
}

func TestWriteJSONSuccess(t *testing.T) {
	rec := httptest.NewRecorder()

	WriteJSON(rec, http.StatusOK, "success", map[string]string{"version": "v1.1"})

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if got := rec.Header().Get("Content-Type"); got != "application/json" {
		t.Fatalf("Content-Type = %q, want %q", got, "application/json")
	}

	var got responseEnvelope
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

func TestWriteJSONErrorPayload(t *testing.T) {
	rec := httptest.NewRecorder()

	WriteJSON(rec, http.StatusBadRequest, "bad request", map[string]string{"field": "invalid"})

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}

	var got responseEnvelope
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if got.Message != "bad request" {
		t.Fatalf("message = %q, want %q", got.Message, "bad request")
	}
	if string(got.Data) != "null" {
		t.Fatalf("data = %s, want null", got.Data)
	}
	if string(got.Error) != `{"field":"invalid"}` {
		t.Fatalf("error = %s, want %s", got.Error, `{"field":"invalid"}`)
	}
}

func TestWriteJSONDefaultsStatusAndMessage(t *testing.T) {
	rec := httptest.NewRecorder()

	WriteJSON(rec, 0, "", nil)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	var got responseEnvelope
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if got.Message != http.StatusText(http.StatusOK) {
		t.Fatalf("message = %q, want %q", got.Message, http.StatusText(http.StatusOK))
	}
	if string(got.Data) != "null" {
		t.Fatalf("data = %s, want null", got.Data)
	}
}
