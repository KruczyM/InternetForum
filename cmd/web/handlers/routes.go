package handlers

import (
	"forum/ui"
	"io/fs"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) Routes() http.Handler {
	//The router automatically handles HTTP methods such as GET, POST, PATCH, and DELETE.
	// This eliminates the need for manual if checks or custom middlewares to verify the request method inside every handler.
	//It provides native support for dynamic routing (e.g., /post/{id}

	router := chi.NewRouter()

	router.NotFound(func(w http.ResponseWriter, r *http.Request) {
        http.Redirect(w, r, "/", http.StatusSeeOther) 
    })

	//Middlewares function like layers of an onion—the request passes through them from the outermost layer to the innermost.
	//Using router.Use(middleware) registers the middleware for all routes within the router.
	router.Use(h.recoverPanic)
	router.Use(h.logRequest)
	router.Use(h.authenticate)

	//  static connect
	staticFS, err := fs.Sub(ui.Files, "static")
	if err != nil {
		panic(err)
	}

	fileServer := http.FileServer(http.FS(staticFS))
	router.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	//live img upload
	uploadServer := http.FileServer(http.Dir("./ui/static/uploads"))
	router.Handle("/static/uploads/*", http.StripPrefix("/static/uploads/", uploadServer))

	//In the Chi router, we map handlers directly to HTTP methods. We use router.Get("/", handler) for retrieval, and similarly Post, Put, Patch, or Delete for other actions.
	router.Get("/", h.home)
	//Create a route group for /post, each sub will start from /post
	router.Route("/post", func(router chi.Router) {
		// .With(middleware) adds middleware to this one path
		router.With(h.requireAuth).Get("/create", h.createPost)
		router.With(h.requireAuth).Post("/create", h.createPost)
		router.Get("/{id}", h.ViewPost)
		router.With(h.requireAuth).Post("/{id}/comment", h.CreateComment)
		router.With(h.requireAuth).Post("/{id}/delete", h.DeletePost)
		router.With(h.requireAuth).Get("/{id}/edit", h.EditPost)
		router.With(h.requireAuth).Post("/{id}/edit", h.UpdatePost)
		router.With(h.requireAuth).Post("/{id}/like", h.PostLike)
		
	})
	router.With(h.requireAuth).Post("/comment/{id}/delete", h.DeleteComment)
	router.With(h.requireAuth).Post("/comment/{id}/like", h.CommentLike)

	router.With(h.requireAuth).Get("/book/create", h.CreateBook)
	router.With(h.requireAuth).Post("/book/create", h.CreateBook)

	router.Route("/auth", func(router chi.Router) {
		router.Get("/register", h.userRegister)
		router.Post("/register", h.userRegisterPost)
		router.Get("/login", h.userLogin)
		router.Post("/login", h.userLoginPost)
		router.Post("/logout", h.userLogoutPost)
	})

	// Chat route (Go-only)
	router.Route("/chat", func(router chi.Router) {
		router.MethodFunc("GET", "/", h.ChatHandler)
		router.MethodFunc("POST", "/", h.ChatHandler)
	})

	//user_panel
	router.Route("/profile", func(r chi.Router) {
	r.Use(h.requireAuth)

	r.Get("/", h.userProfile)
	r.Post("/edit", h.userProfileEditPost)
	r.Post("/password", h.userProfilePasswordPost)
})

	return router
}
