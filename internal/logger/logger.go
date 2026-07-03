package logger

import (
	"encoding/json"
	"log"
)

type Logger struct {
	*log.Logger
}

type Access struct {
	Level         string
	Method        string
	Path          string
	RequestBody   string `json:"RequestBody,omitempty"`
	RequestParams string `json:"RequestParams,omitempty"`
}

type Info struct {
	Level         string
	Method        string
	Path          string
	Status        string
	Duration_ms   string
	RequestBody   string `json:"RequestBody,omitempty"`
	RequestParams string `json:"RequestParams,omitempty"`
	ResponseBody  string
}

type Error struct {
	Level         string
	ErrorMessage  string
	StackTrace    string `json:"StackTrace,omitempty"`
}

func NewLogger() *Logger {
	return &Logger{
		Logger: log.New(log.Writer(), "[sop-in-go]", log.LstdFlags),
	}
}

func (l *Logger) Debug(msg string, args any) {
	log, err := json.Marshal(args)
	if err != nil {
		l.Printf("[DEBUG] error marshalling args in Debug function: %s", err.Error())
	}
	l.Printf("[DEBUG] %s:%s", msg, string(log))
}

func (l *Logger) Access(a Access) {
	log, err := json.Marshal(a)
	if err != nil {
		l.Printf("[ACCESS] error marshalling access in Access function: %s", err.Error())
	}
	l.Printf("[ACCESS] %s", string(log))
}

func (l *Logger) Info(i Info) {
	log, err := json.Marshal(i)
	if err != nil {
		l.Printf("[INFO] error marshalling info in Info function: %s", err.Error())
	}
	l.Printf("[INFO] %s", string(log))
}

func (l *Logger) Error(msg Error) {
	log, err := json.Marshal(msg)
	if err != nil {
		l.Printf("[ERROR] error marshalling error in Error function: %s", err.Error())
	}
	l.Printf("[ERROR] %s", string(log))
}
