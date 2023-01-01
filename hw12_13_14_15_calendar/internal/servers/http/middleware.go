package serverhttp

import (
	"fmt"
	"net/http"
	"strings"
)

func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := &responseWriter{w, http.StatusOK}
		next.ServeHTTP(rw, r)

		s.logger.Info(
			fmt.Sprintf("%s %s %s %s %d %s",
				requestAddr(r),
				r.Method,
				r.RequestURI,
				r.Proto,
				rw.code,
				userAgent(r),
			))
	})
}

func requestAddr(r *http.Request) string {
	return strings.Split(r.RemoteAddr, ":")[0]
}

func userAgent(r *http.Request) string {
	userAgents := r.Header["User-Agent"]
	if len(userAgents) > 0 {
		return "\"" + userAgents[0] + "\""
	}
	return ""
}
