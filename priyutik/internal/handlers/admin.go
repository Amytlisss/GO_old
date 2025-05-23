package handlers

import (
	"net/http"
	"time"

	"github.com/Amytlisss/GO_old/priyutik/internal/models"
)

func (h *Handlers) AdminPage(w http.ResponseWriter, r *http.Request) {
	user, ok := h.GetUserFromSession(w, r)
	if !ok || user.Role != "admin" {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	dateFilter := r.URL.Query().Get("date")

	var meetings []models.Meeting
	var err error

	if dateFilter != "" {
		date, err := time.Parse("2006-01-02", dateFilter)
		if err != nil {
			h.ErrorResponse(w, "Неверный формат даты", http.StatusBadRequest)
			return
		}
		meetings, err = h.repo.GetMeetingsByDate(date)
	} else {
		meetings, err = h.repo.GetAllMeetings()
	}
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
		DateFilter string
	}{
		User:       user,
		Meetings:   meetings,
		Animals:    animals,
		DateFilter: dateFilter,
	}

	h.RenderTemplate(w, "admin_page.html", data)
}
