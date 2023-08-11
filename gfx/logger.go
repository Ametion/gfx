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
		color := "\033[33m" // Default to yellow
		if l.statusCode == 200 || l.statusCode == 201 {
			color = "\033[32m" // Green for success
		} else if l.statusCode == 500 || l.statusCode == 400 || l.statusCode == 401 || l.statusCode == 402 || l.statusCode == 403 || l.statusCode == 404 {
			color = "\033[31m" // Red for errors
		}

		// Bold text
		bold := "\033[1m"

		// Reset formatting after printing
		reset := "\033[0m"

		fmt.Printf("%sDate: %s, Method: %s, Status code: %s%d%s, Time taken: %v%s\n", bold, time.Now().Format(time.RFC1123), l.method, color, l.statusCode, reset, reset)
	}

	return n, err
}
