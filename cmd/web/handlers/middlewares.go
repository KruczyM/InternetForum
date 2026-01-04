package handlers

import (
	"context"
	"fmt"
	"forum/internal/models"
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
		// If the user is not authenticated, redirect them to the login page and return from the middleware chain so that no subsequent handlers in the chain are executed.
		if !h.isAuthenticated(r){
			h.SessionManager.Put(r.Context(), "flash", &FlashMessage{
				Type: "error",
				Msg:  "Please log in to perform this action.",
			})
			http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
			return
		}
		// Otherwise set the "Cache-Control: no-store" header so that pages require authentication are not stored in the users browser cache (or other intermediary cache).
		w.Header().Add("Cache-Control", "no-store")
		// just call next handler 
		next.ServeHTTP(w, r)
	})
}

func (h *Handler) authenticate (next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Retrieve the authenticatedUserID value from the session using the GetString() method. This will return the "" value for an string ("") ifno
		// "authenticatedUserID" value is in the session -- in which case we call the next handler in the chain as normal and return.
		id := h.SessionManager.GetString(r.Context(), "authenticatedUserID")
		if id == "" {
			next.ServeHTTP(w, r)
			return 
		}
		// Otherwise, we check to see if a user with that ID exists in our database.
		exists, err := models.ExistsUser(h.DB, id)
		if err != nil {
			h.serverError(w,err)
			return
		}
		// If a matching user is found, we know that the request is coming from an authenticated user who exists in our database. We create a new copy of the request (with an
		//isAuthenticatedContextKey value of true in the request context) and assign it to r.
		if exists{
			ctx := context.WithValue(r.Context(), "isAuthenticated", true)
			r = r.WithContext(ctx)
		}
		next.ServeHTTP(w, r)
	})
}