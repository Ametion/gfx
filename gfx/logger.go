package gfx

import (
	"fmt"
	"net/http"
	"time"
)

type LoggingResponseWriter struct {
	http.ResponseWriter
	g          *GFXEngine
	statusCode int
	method     string
}

func (l *LoggingResponseWriter) WriteHeader(statusCode int) {
	l.statusCode = statusCode
	l.ResponseWriter.WriteHeader(statusCode)
}

func (l *LoggingResponseWriter) Write(b []byte) (int, error) {
	n, err := l.ResponseWriter.Write(b)
	if err == nil && l.g.development {
		fmt.Printf("Date: %s, Method: %s, Status code: %d\n", time.Now().Format(time.RFC1123), l.method, l.statusCode)
	}
	return n, err
}
