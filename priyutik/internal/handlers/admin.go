package handlers

import (
	"net/http"
	"time"

	"priyutik/internal/models"
)

func (h *Handlers) AdminPage(w http.ResponseWriter, r *http.Request) {
	user, ok := h.GetUserFromSession(w, r)
	if !ok || user.Role != "admin" {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	dateFilter := r.URL.Query().Get("date")

	meetings, err := h.repo.GetFilteredMeetings(dateFilter)
	if err != nil {
		h.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	animals, err := h.repo.GetAllAnimals()
	if err != nil {
		h.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	data := struct {
		User       models.User
		Meetings   []models.Meeting
		Animals    []models.Animal
		Now        time.Time
		DateFilter string
	}{
		User:       user,
		Meetings:   meetings,
		Animals:    animals,
		Now:        time.Now(),
		DateFilter: dateFilter,
	}

	h.RenderTemplate(w, "admin_page.html", data)
}
