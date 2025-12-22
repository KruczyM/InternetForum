package handlers

import (
	"net/http"
	"github.com/go-chi/chi/v5"
)

func (h *Handler) Routes() http.Handler {
//The router automatically handles HTTP methods such as GET, POST, PATCH, and DELETE. 
// This eliminates the need for manual if checks or custom middlewares to verify the request method inside every handler.
//It provides native support for dynamic routing (e.g., /post/{id}

	router := chi.NewRouter()

//Middlewares function like layers of an onion—the request passes through them from the outermost layer to the innermost.
//Using router.Use(middleware) registers the middleware for all routes within the router.
	router.Use(h.recoverPanic)
	router.Use(h.logRequest)
	router.Use(h.authenticate)

//In the Chi router, we map handlers directly to HTTP methods. We use router.Get("/", handler) for retrieval, and similarly Post, Put, Patch, or Delete for other actions.
	router.Get("/", h.home)
//Create a route group for /post, each sub will start from /post
	router.Route("/post", func(router chi.Router) {
		// .With(middleware) adds middleware to this one path
		router.With(h.requireAuth).Get("/create", h.createPost)
		router.With(h.requireAuth).Post("/create", h.createPost)
		router.Get("/{id}", h.ViewPost)
		router.With(h.requireAuth).Post("/{id}/comment", h.CreateComment)
	})

	router.Route("/auth", func(router chi.Router) {
		router.Get("/register", h.userRegister)
		router.Post("/register", h.userRegisterPost)
		router.Get("/login", h.userLogin)
		router.Post("/login", h.userLoginPost)
		router.Post("/logout", h.userLogoutPost)
	})

		// Chat routes
	router.Route("/chat", func(router chi.Router) {
		router.Post("/message", h.PostChatMessage)
		router.Get("/messages", h.GetChatMessages)
		router.Get("/categories", h.GetChatCategories)
	})


	return router
}
