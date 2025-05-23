package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/Amytlisss/GO_old/priyutik/internal/models"
)

func (h *Handlers) MeetingsPage(w http.ResponseWriter, r *http.Request) {
	user, ok := h.GetUserFromSession(w, r)
	if !ok {
		return
	}

	if r.Method == http.MethodPost {
		dateStr := r.FormValue("date")
		timeStr := r.FormValue("time")

		if dateStr == "" || timeStr == "" {
			h.ErrorResponse(w, "Дата и время должны быть указаны", http.StatusBadRequest)
			return
		}

		dateTime, err := time.Parse("2006-01-02T15:04", dateStr+"T"+timeStr)
		if err != nil {
			h.ErrorResponse(w, "Неверный формат даты или времени", http.StatusBadRequest)
			return
		}

		if err := h.repo.CreateMeeting(user.ID, dateTime); err != nil {
			h.ErrorResponse(w, "Ошибка при создании встречи", http.StatusInternalServerError)
			return
		}
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
		Role     string
	}{
		User:     user,
		Meetings: meetings,
		Now:      time.Now(),
		Role:     user.Role,
	}

	h.RenderTemplate(w, "meetings.html", data)
}

func (h *Handlers) CancelMeeting(w http.ResponseWriter, r *http.Request) {
	_, ok := h.GetUserFromSession(w, r)
	if !ok {
		return
	}

	meetingID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		h.ErrorResponse(w, "Неверный ID встречи", http.StatusBadRequest)
		return
	}

	if err := h.repo.CancelMeeting(meetingID); err != nil {
		h.ErrorResponse(w, "Ошибка при отмене встречи", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/meetings", http.StatusFound)
}

func (h *Handlers) EditMeetingPage(w http.ResponseWriter, r *http.Request) {
	user, ok := h.GetUserFromSession(w, r)
	if !ok {
		return
	}

	meetingID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		h.ErrorResponse(w, "Неверный ID встречи", http.StatusBadRequest)
		return
	}

	meeting, err := h.repo.GetMeetingByID(meetingID)
	if err != nil {
		h.ErrorResponse(w, "Ошибка при получении встречи", http.StatusInternalServerError)
		return
	}

	if meeting.UserID != user.ID && user.Role != "admin" {
		h.ErrorResponse(w, "У вас нет прав для редактирования этой встречи", http.StatusForbidden)
		return
	}

	data := struct {
		User    models.User
		Meeting models.Meeting
	}{
		User:    user,
		Meeting: *meeting,
	}

	h.RenderTemplate(w, "edit_meeting.html", data)
}

func (h *Handlers) EditMeeting(w http.ResponseWriter, r *http.Request) {
	user, ok := h.GetUserFromSession(w, r)
	if !ok {
		return
	}

	meetingID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		h.ErrorResponse(w, "Неверный ID встречи", http.StatusBadRequest)
		return
	}

	meeting, err := h.repo.GetMeetingByID(meetingID)
	if err != nil {
		h.ErrorResponse(w, "Ошибка при получении встречи", http.StatusInternalServerError)
		return
	}

	if meeting.UserID != user.ID && user.Role != "admin" {
		h.ErrorResponse(w, "У вас нет прав для редактирования этой встречи", http.StatusForbidden)
		return
	}

	if r.Method == http.MethodPost {
		newDateStr := r.FormValue("date")
		newTimeStr := r.FormValue("time")

		if newDateStr == "" || newTimeStr == "" {
			h.ErrorResponse(w, "Дата и время не могут быть пустыми", http.StatusBadRequest)
			return
		}

		newDateTime, err := time.Parse("2006-01-02 15:04", newDateStr+" "+newTimeStr)
		if err != nil {
			h.ErrorResponse(w, "Неверный формат времени", http.StatusBadRequest)
			return
		}

		if err := h.repo.UpdateMeeting(meetingID, newDateTime); err != nil {
			h.ErrorResponse(w, "Ошибка при изменении времени встречи", http.StatusInternalServerError)
			return
		}

		if user.Role == "admin" {
			http.Redirect(w, r, "/admin", http.StatusFound)
		} else {
			http.Redirect(w, r, "/meetings", http.StatusFound)
		}
		return
	}

	data := struct {
		User    models.User
		Meeting models.Meeting
	}{
		User:    user,
		Meeting: *meeting,
	}

	h.RenderTemplate(w, "edit_meeting.html", data)
}
