package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime/debug"
	"strconv"
	"time"

	"github.com/DilankaHer/sop-in-go/internal/logger"
)

type Middleware interface {
	Recover(next http.Handler) http.Handler
	AccessLog(next http.Handler) http.Handler
	Response(h func(r *http.Request) StandardResponse) http.HandlerFunc
}

type middleware struct {
	logger *logger.Logger
}

type StandardResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data"`
	Error   any    `json:"error"`
}

func NewMiddleware(logger *logger.Logger) Middleware {
	return &middleware{logger: logger}
}

func (m *middleware) Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				m.logger.Error(logger.Error{
					Level:        "error",
					ErrorMessage: fmt.Sprintf("panic: %v", rec),
					StackTrace:   string(debug.Stack()),
				})
				resp := StandardResponse{
					Status:  http.StatusInternalServerError,
					Message: "internal server error",
					Error:   nil,
				}
				body, err := json.Marshal(resp)
				if err != nil {
					m.logger.Error(logger.Error{
						Level:        "error",
						ErrorMessage: fmt.Sprintf("error marshalling response in Recover middleware: %s", err.Error()),
						StackTrace:   string(debug.Stack()),
					})
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(resp.Status)
				_, err = w.Write(body)
				if err != nil {
					m.logger.Error(logger.Error{
						Level:        "error",
						ErrorMessage: fmt.Sprintf("error writing response in Recover middleware: %s", err.Error()),
						StackTrace:   string(debug.Stack()),
					})
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (m *middleware) AccessLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		r = r.WithContext(context.WithValue(r.Context(), "startTime", start))
		if len(r.URL.Query()) > 0 {
			p, err := json.Marshal(r.URL.Query())
			if err != nil {
				m.logger.Error(logger.Error{
					Level:        "error",
					ErrorMessage: fmt.Sprintf("error marshaling request params: %s", err.Error()),
					StackTrace:   string(debug.Stack()),
				})
			}
			access := logger.Access{
				Level:         "Access",
				Method:        r.Method,
				Path:          r.URL.Path,
				RequestParams: string(p),
			}
			r.SetPathValue("params", string(p))
			m.logger.Access(access)
		} else {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				m.logger.Error(logger.Error{
					Level:        "error",
					ErrorMessage: fmt.Sprintf("error reading request body: %s", err.Error()),
					StackTrace:   string(debug.Stack()),
				})
			}
			access := logger.Access{
				Level:       "Access",
				Method:      r.Method,
				Path:        r.URL.Path,
				RequestBody: string(body),
			}
			r.SetPathValue("body", string(body))
			m.logger.Access(access)
			r.Body = io.NopCloser(bytes.NewBuffer(body))
		}
		next.ServeHTTP(w, r)
	})
}

func (m *middleware) Response(h func(r *http.Request) StandardResponse) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := h(r)
		if resp.Status == 0 {
			resp.Status = http.StatusOK
		}
		if resp.Message == "" {
			resp.Message = http.StatusText(resp.Status)
		}

		if err := m.validateResponse(resp); err != nil {
			m.logger.Error(logger.Error{
				Level:        "error",
				ErrorMessage: fmt.Sprintf("error validating response: %s", err.Error()),
				StackTrace:   string(debug.Stack()),
			})
			if os.Getenv("ENV") == "local" || os.Getenv("ENV") == "" {
				m.logger.Fatalf("%s", err.Error())
			}
			resp = StandardResponse{
				Status:  http.StatusInternalServerError,
				Message: "internal server error",
				Error:   nil,
			}
		}

		body, err := json.Marshal(resp)
		if err != nil {
			m.logger.Error(logger.Error{
				Level:        "error",
				ErrorMessage: fmt.Sprintf("error marshalling response: %s", err.Error()),
				StackTrace:   string(debug.Stack()),
			})
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resp.Status)
		_, err = w.Write(body)
		if err != nil {
			m.logger.Error(logger.Error{
				Level:        "error",
				ErrorMessage: fmt.Sprintf("error writing response: %s", err.Error()),
				StackTrace:   string(debug.Stack()),
			})
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		info := logger.Info{
			Level:         "info",
			Method:        r.Method,
			Path:          r.URL.Path,
			Status:        strconv.Itoa(resp.Status),
			Duration_ms:   strconv.Itoa(int(requestDuration(r).Milliseconds())),
			RequestBody:   r.PathValue("body"),
			RequestParams: r.PathValue("params"),
			ResponseBody:  string(body),
		}
		m.logger.Info(info)
	})
}

func requestDuration(r *http.Request) time.Duration {
	start, ok := r.Context().Value("startTime").(time.Time)
	if !ok {
		return 0
	}
	return time.Since(start)
}

func (m *middleware) validateResponse(resp StandardResponse) error {
	if resp.Status >= 400 && resp.Data != nil {
		return fmt.Errorf("status %d must not include data", resp.Status)
	}
	if resp.Status < 400 && resp.Error != nil {
		return fmt.Errorf("status %d must not include error", resp.Status)
	}
	return nil
}
