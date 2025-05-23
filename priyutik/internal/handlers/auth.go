package handlers

import (
	"html/template"
	"net/http"
	"priyutik/internal/repository"
	"time"

	"priyutik/internal/models"

	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	repo  *repository.Repository
	store *sessions.CookieStore
}

func (h *Handlers) RegisterPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl, err := template.ParseFiles("templates/register.html")
		if err != nil {
			h.ErrorResponse(w, err, http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
	} else if r.Method == http.MethodPost {
		r.ParseForm()

		user := models.User{
			Name:     r.FormValue("name"),
			Phone:    r.FormValue("phone"),
			Email:    r.FormValue("email"),
			Password: r.FormValue("password"),
			Role:     "user",
		}

		if user.Name == "" || user.Phone == "" || user.Email == "" || user.Password == "" {
			h.ErrorResponse(w, "Все поля должны быть заполнены", http.StatusBadRequest)
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			h.ErrorResponse(w, err, http.StatusInternalServerError)
			return
		}
		user.Password = string(hashedPassword)

		if err := h.repo.CreateUser(&user); err != nil {
			h.ErrorResponse(w, err, http.StatusInternalServerError)
			return
		}

		session, _ := h.store.Get(r, "session-name")
		session.Values["user"] = user
		session.Save(r, w)

		http.Redirect(w, r, "/user_profile", http.StatusFound)
	}
}

func (h *Handlers) LoginPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl, err := template.ParseFiles("templates/login.html")
		if err != nil {
			h.ErrorResponse(w, err, http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
	} else if r.Method == http.MethodPost {
		r.ParseForm()
		phone := r.FormValue("phone")
		password := r.FormValue("password")

		if phone == "" || password == "" {
			h.ErrorResponse(w, "Номер телефона и пароль не могут быть пустыми", http.StatusBadRequest)
			return
		}

		user, err := h.repo.GetUserByPhone(phone)
		if err != nil {
			h.ErrorResponse(w, "Неверный номер телефона или пароль", http.StatusUnauthorized)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
		if err != nil {
			h.ErrorResponse(w, "Неверный номер телефона или пароль", http.StatusUnauthorized)
			return
		}

		session, _ := h.store.Get(r, "session-name")
		session.Values["user"] = user
		session.Save(r, w)

		if user.Role == "admin" {
			http.Redirect(w, r, "/admin", http.StatusFound)
		} else {
			http.Redirect(w, r, "/user_profile", http.StatusFound)
		}
	}
}

func (h *Handlers) Logout(w http.ResponseWriter, r *http.Request) {
	session, _ := h.store.Get(r, "session-name")
	delete(session.Values, "user")
	session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusFound)
}
func (h *Handlers) HomePage(w http.ResponseWriter, r *http.Request) {
	session, _ := h.store.Get(r, "session-name")
	user, ok := session.Values["user"].(models.User)

	animals, err := h.repo.GetAllAnimals()
	if err != nil {
		h.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	data := struct {
		User       models.User
		IsLoggedIn bool
		Role       string
		Animals    []models.Animal
	}{
		User:       user,
		IsLoggedIn: ok,
		Role:       user.Role,
		Animals:    animals,
	}

	h.RenderTemplate(w, "home_page.html", data)
}

func (h *Handlers) UserProfile(w http.ResponseWriter, r *http.Request) {
	session, _ := h.store.Get(r, "session-name")
	user, ok := session.Values["user"].(models.User)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	meetings, err := h.repo.GetMeetings(user.ID)
	if err != nil {
		h.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	data := struct {
		User     models.User
		Meetings []models.Meeting
		Now      time.Time
	}{
		User:     user,
		Meetings: meetings,
		Now:      time.Now(),
	}

	h.RenderTemplate(w, "user_profile.html", data)
}
