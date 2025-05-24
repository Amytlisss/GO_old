package handlers

import (
	"encoding/gob"
	"html/template"
	"net/http"

	"priyutik/internal/config"
	"priyutik/internal/models"

	"priyutik/internal/repository"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

type Handlers struct {
	repo  *repository.Repository
	store *sessions.CookieStore
	cfg   *config.Config
}

func NewHandlers(repo *repository.Repository, store *sessions.CookieStore, cfg *config.Config) *Handlers {
	gob.Register(models.User{})
	return &Handlers{
		repo:  repo,
		store: store,
		cfg:   cfg,
	}
}

func (h *Handlers) RegisterRoutes() {
	r := mux.NewRouter()

	// Public routes
	r.HandleFunc("/", h.HomePage).Methods("GET")
	r.HandleFunc("/register", h.RegisterPage).Methods("GET", "POST")
	r.HandleFunc("/login", h.LoginPage).Methods("GET", "POST")
	r.HandleFunc("/animals", h.AnimalsPage).Methods("GET")

	// Protected routes
	r.Handle("/logout", h.Authenticate(http.HandlerFunc(h.Logout))).Methods("GET")
	r.Handle("/user_profile", h.Authenticate(http.HandlerFunc(h.UserProfile))).Methods("GET")
	r.Handle("/meetings", h.Authenticate(http.HandlerFunc(h.MeetingsPage))).Methods("GET", "POST")
	r.Handle("/cancel_meeting", h.Authenticate(http.HandlerFunc(h.CancelMeeting))).Methods("GET", "POST")
	r.Handle("/meetings/edit", h.Authenticate(http.HandlerFunc(h.EditMeetingPage))).Methods("GET")
	r.Handle("/meetings/edit", h.Authenticate(http.HandlerFunc(h.EditMeeting))).Methods("POST")
	r.Handle("/admin", h.Authenticate(h.AdminOnly(http.HandlerFunc(h.AdminPage)))).Methods("GET")
	r.Handle("/animals/edit", h.Authenticate(h.AdminOnly(http.HandlerFunc(h.EditAnimalPage)))).Methods("GET")
	r.Handle("/animals/edit", h.Authenticate(h.AdminOnly(http.HandlerFunc(h.EditAnimal)))).Methods("POST")
	r.Handle("/animals/delete", h.Authenticate(h.AdminOnly(http.HandlerFunc(h.DeleteAnimal)))).Methods("GET", "POST")
	r.Handle("/admin/animals/add", h.Authenticate(h.AdminOnly(http.HandlerFunc(h.AddAnimal)))).Methods("GET", "POST")
	http.Handle("/", r)
}

func (h *Handlers) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := h.store.Get(r, "session-name")
		if session.Values["user"] == nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (h *Handlers) AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := h.store.Get(r, "session-name")
		user, ok := session.Values["user"].(models.User)
		if !ok || user.Role != "admin" {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (h *Handlers) GetUserFromSession(w http.ResponseWriter, r *http.Request) (models.User, bool) {
	session, _ := h.store.Get(r, "session-name")
	user, ok := session.Values["user"].(models.User)
	return user, ok
}

func (h *Handlers) ErrorResponse(w http.ResponseWriter, err interface{}, status int) {
	http.Error(w, err.(error).Error(), status)
}

func (h *Handlers) RenderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	t, err := template.ParseFiles("templates/" + tmpl)
	if err != nil {
		h.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}
	if err := t.Execute(w, data); err != nil {
		h.ErrorResponse(w, err, http.StatusInternalServerError)
	}
}
