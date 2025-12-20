package handlers

import (
	"fmt"
	"net/http"
	"time"
)

// recoverPanic recovers from panics and logs them.
func (h *Handler) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")

				h.serverError(w, fmt.Errorf("%v", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}
//to add status code to response
type responseRecorder struct {
	http.ResponseWriter
	status int
}

func (rr *responseRecorder) WriteHeader(code int) {
	rr.status = code
	rr.ResponseWriter.WriteHeader(code)
}

// logRequest logs the start and end of each request.
func (h *Handler) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		rr := &responseRecorder{
			ResponseWriter: w,
			status: http.StatusOK,
		}

		start := time.Now()
		next.ServeHTTP(rr, r)

		h.InfoLog.Printf(
			"%d %s %s %s",
			rr.status,
			r.Method,
			r.RequestURI,
			time.Since(start),
		)
	})
}

func (h *Handler) requireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// session, err := h.SessionManager.Get(r, "session")
		// if err != nil {
		// 	h.serverError(w, err)
		// 	return
		// }
		// if session == nil {
		// 	http.Redirect(w, r, "/auth/register", http.StatusSeeOther)
		// 	return
		// }
		next.ServeHTTP(w, r)
	})
}