package core

import (
	"encoding/json"
	"log/slog"
	"runtime"
	"strconv"
)

const reservedKeyCount = 4

// SourceLocation represents the caller source location in GCP Cloud Logging.
type SourceLocation struct {
	File     string `json:"file,omitempty"`
	Line     string `json:"line,omitempty"`
	Function string `json:"function,omitempty"`
}

func collectRecordAttrs(record slog.Record) []slog.Attr {
	attrs := make([]slog.Attr, 0, record.NumAttrs())
	record.Attrs(func(a slog.Attr) bool {
		attrs = append(attrs, a)
		return true
	})
	return attrs
}

func addSourceLocation(m map[string]any, pc uintptr) {
	fs := runtime.CallersFrames([]uintptr{pc})
	f, _ := fs.Next()
	if f.File != "" {
		m["logging.googleapis.com/sourceLocation"] = SourceLocation{
			File:     f.File,
			Line:     strconv.Itoa(f.Line),
			Function: f.Function,
		}
	}
}

// GCPLogEntry represents the top-level structure written to JSON output.
type GCPLogEntry struct {
	HTTPRequest *HttpRequest   `json:"httpRequest,omitempty"`
	Payload     map[string]any `json:"-"`
	Message     any            `json:"message,omitempty"`
	Timestamp   string         `json:"timestamp,omitempty"`
	Severity    string         `json:"severity"`
}

func populateOptionalFields(m map[string]any, e GCPLogEntry) {
	if e.HTTPRequest != nil {
		m["httpRequest"] = e.HTTPRequest
	}
	if e.Timestamp != "" {
		m["timestamp"] = e.Timestamp
	}
}

func populateEntryFields(m map[string]any, e GCPLogEntry) {
	if e.Message != nil {
		m["message"] = e.Message
	}
	populateOptionalFields(m, e)
}

// MarshalJSON serializes the GCPLogEntry into a JSON object including custom payload keys.
func (e GCPLogEntry) MarshalJSON() ([]byte, error) {
	m := make(map[string]any, reservedKeyCount+len(e.Payload))
	m["severity"] = e.Severity
	populateEntryFields(m, e)
	mergePayload(m, e.Payload)
	return json.Marshal(m)
}
