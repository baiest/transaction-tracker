package services

import (
	"context"
	"encoding/json"
	"os"
	"sync"
	"time"
	"transaction-tracker/logger/models"
)

type JSONLogger struct {
	ServiceName string
	file        *os.File
	mu          sync.Mutex
}

type LogEntry struct {
	Timestamp  string                 `json:"timestamp"`
	Level      string                 `json:"level"`
	Service    string                 `json:"service"`
	Event      string                 `json:"event"`
	Error      string                 `json:"message,omitempty"`
	Properties map[string]interface{} `json:"properties,omitempty"`
}

func NewJSONLogger(_ context.Context, serviceName, filePath string) (*JSONLogger, error) {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	return &JSONLogger{
		ServiceName: serviceName,
		file:        file,
	}, nil
}

func (l *JSONLogger) Log(level string, props models.LogProperties) {
	l.mu.Lock()
	defer l.mu.Unlock()

	entry := LogEntry{
		Timestamp:  time.Now().UTC().Format(time.RFC3339),
		Level:      level,
		Service:    l.ServiceName,
		Event:      props.Event,
		Properties: map[string]interface{}{},
	}

	if (level == "error" || level == "panic") && props.Error != nil {
		entry.Error = props.Error.Error()
	}

	for _, additionalParam := range props.AdditionalParams {
		if additionalParam == nil {
			continue
		}

		for key, value := range additionalParam.LogProperties() {
			entry.Properties[key] = value
		}
	}

	data, err := json.Marshal(entry)
	if err != nil {
		data = []byte(`{"error":"failed to marshal log entry"}`)
	}

	l.file.WriteString(string(data) + "\n")
}

func (l *JSONLogger) Close() error {
	return l.file.Close()
}
