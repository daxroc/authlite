package logging

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"sync"
)

type Format string

const (
	LogFmt Format = "logfmt"
	JSON   Format = "json"
)

var (
	logFormat Format = LogFmt // default
	mu        sync.RWMutex
)

// SetFormat allows setting log format externally (e.g., from config)
func SetFormat(f Format) {
	mu.Lock()
	defer mu.Unlock()
	logFormat = f
}

func Info(msg string, kv ...any) {
	logWithCaller(os.Stdout, "info", msg, kv...)
}

func Error(msg string, kv ...any) {
	logWithCaller(os.Stderr, "error", msg, kv...)
}

func logWithCaller(out io.Writer, level string, msg string, kv ...any) {
	_, file, line, _ := runtime.Caller(2)
	caller := fmt.Sprintf("%s:%d", trimPath(file), line)

	mu.RLock()
	defer mu.RUnlock()

	if logFormat == JSON {
		logAsJSON(out, level, msg, caller, kv...)
	} else {
		logAsLogFmt(out, level, msg, caller, kv...)
	}
}

func logAsLogFmt(w io.Writer, level, msg, caller string, kv ...any) {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("level=%s msg=%q caller=%s", level, msg, caller))
	for i := 0; i < len(kv)-1; i += 2 {
		k, ok := kv[i].(string)
		if !ok {
			continue
		}
		b.WriteString(fmt.Sprintf(" %s=%s", k, formatValue(kv[i+1])))
	}
	b.WriteString("\n")
	io.WriteString(w, b.String())
}

func logAsJSON(w io.Writer, level, msg, caller string, kv ...any) {
	m := map[string]interface{}{
		"level":  level,
		"msg":    msg,
		"caller": caller,
	}
	for i := 0; i < len(kv)-1; i += 2 {
		k, ok := kv[i].(string)
		if !ok {
			continue
		}
		m[k] = kv[i+1]
	}
	json.NewEncoder(w).Encode(m)
}

func formatValue(v any) string {
	switch x := v.(type) {
	case string:
		return `"` + x + `"`
	default:
		return fmt.Sprintf("%v", v)
	}
}

func trimPath(path string) string {
	if idx := strings.LastIndex(path, "/"); idx != -1 {
		return path[idx+1:]
	}
	return path
}
