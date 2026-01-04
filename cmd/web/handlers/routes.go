package handlers

import (
	"forum/ui"
	"io/fs"
	"net/http"
	"strings"
)

func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()

	staticFS, err := fs.Sub(ui.Files, "static")
	if err != nil {
		panic(err)
	}

	mux.Handle(
		"/static/",
		http.StripPrefix(
			"/static/",
			http.FileServer(http.FS(staticFS)),
		),
	)

	mux.Handle(
		"/static/uploads/",
		http.StripPrefix(
			"/static/uploads/",
			http.FileServer(http.Dir("./ui/static/uploads")),
		),
	)

	mux.HandleFunc("/", h.home)


	mux.HandleFunc("/auth/register", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			h.userRegister(w, r)
		case http.MethodPost:
			h.userRegisterPost(w, r)
		default:
			http.NotFound(w, r)
		}
	})

	mux.HandleFunc("/auth/login", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			h.userLogin(w, r)
		case http.MethodPost:
			h.userLoginPost(w, r)
		default:
			http.NotFound(w, r)
		}
	})

	mux.HandleFunc("/auth/logout", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			h.userLogoutPost(w, r)
			return
		}
		http.NotFound(w, r)
	})

	mux.Handle(
		"/post/create",
		h.requireAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodGet || r.Method == http.MethodPost {
				h.createPost(w, r)
				return
			}
			http.NotFound(w, r)
		})),
	)

	mux.Handle(
		"/post/",
		h.requireAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			path := strings.TrimPrefix(r.URL.Path, "/post/")
			parts := strings.Split(path, "/")

			if len(parts) == 1 && r.Method == http.MethodGet {
				h.ViewPost(w, r)
				return
			}

			if len(parts) == 2 {
				switch parts[1] {
				case "edit":
					if r.Method == http.MethodGet {
						h.EditPost(w, r)
						return
					}
					if r.Method == http.MethodPost {
						h.UpdatePost(w, r)
						return
					}
				case "delete":
					if r.Method == http.MethodPost {
						h.DeletePost(w, r)
						return
					}
				case "comment":
					if r.Method == http.MethodPost {
						h.CreateComment(w, r)
						return
					}
				case "like":
					if r.Method == http.MethodPost {
						h.PostLike(w, r)
						return
					}
				case "dislike":
					if r.Method == http.MethodPost {
						h.PostDislike(w, r)
						return
					}
				}
			}

			http.NotFound(w, r)
		})),
	)

	mux.Handle(
		"/comment/",
		h.requireAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			path := strings.TrimPrefix(r.URL.Path, "/comment/")
			parts := strings.Split(path, "/")

			if len(parts) != 2 {
				http.NotFound(w, r)
				return
			}

			switch parts[1] {
			case "delete":
				if r.Method == http.MethodPost {
					h.DeleteComment(w, r)
					return
				}
			case "like":
				if r.Method == http.MethodPost {
					h.CommentLike(w, r)
					return
				}
			case "dislike":
				if r.Method == http.MethodPost {
					h.CommentDislike(w, r)
					return
				}
			}

			http.NotFound(w, r)
		})),
	)

	mux.Handle(
		"/book/create",
		h.requireAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodGet || r.Method == http.MethodPost {
				h.CreateBook(w, r)
				return
			}
			http.NotFound(w, r)
		})),
	)

	mux.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/chat/", http.StatusMovedPermanently)
	})

	mux.HandleFunc("/chat/", h.ChatHandler)

	mux.Handle(
		"/profile/",
		h.requireAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch {
			case r.URL.Path == "/profile/" && r.Method == http.MethodGet:
				h.userProfile(w, r)
			case r.URL.Path == "/profile/edit" && r.Method == http.MethodPost:
				h.userProfileEditPost(w, r)
			case r.URL.Path == "/profile/password" && r.Method == http.MethodPost:
				h.userProfilePasswordPost(w, r)
			case r.URL.Path == "/profile/avatar" && r.Method == http.MethodPost:
				h.changeAvatar(w, r)
			default:
				http.NotFound(w, r)
			}
		})),
	)

	mux.Handle(
		"/u/",
		h.requireAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h.publicUserProfile(w, r)
		})),
	)

	mux.HandleFunc("/test500", h.test500)

	var handler http.Handler = mux
	handler = h.authenticate(handler)
	handler = h.logRequest(handler)
	handler = h.recoverPanic(handler)

	return handler
}
