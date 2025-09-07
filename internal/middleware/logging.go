package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/MananLed/evalProjectRailway/pkg/logger"
)

func LoggingMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		logger.LogToFile(fmt.Sprintf("Startted %v %v from %v", r.Method, r.URL.Path, r.RemoteAddr))

		next.ServeHTTP(w, r)

		logger.LogToFile(fmt.Sprintf("Completed %v %v in %v", r.Method, r.URL.Path, time.Since(start)))
	})

}
