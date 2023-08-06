package gfx

import "net/http"

type statusRecorder struct {
	http.ResponseWriter
	statusCode *int
}

func newStatusRecorder(w http.ResponseWriter, statusCode *int) *statusRecorder {
	return &statusRecorder{w, statusCode}
}

func (s *statusRecorder) WriteHeader(code int) {
	*s.statusCode = code
	s.ResponseWriter.WriteHeader(code)
}
