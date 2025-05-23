package handlers

import (
	"net/http"
	"strconv"

	"github.com/Amytlisss/GO_old/priyutik/internal/models"
)

func (h *Handlers) AnimalsPage(w http.ResponseWriter, r *http.Request) {
	animals, err := h.repo.GetAllAnimals()
	if err != nil {
		h.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	user, _ := h.GetUserFromSession(w, r)

	data := struct {
		User    models.User
		Animals []models.Animal
	}{
		User:    user,
		Animals: animals,
	}

	h.RenderTemplate(w, "animals.html", data)
}

func (h *Handlers) EditAnimalPage(w http.ResponseWriter, r *http.Request) {
	user, ok := h.GetUserFromSession(w, r)
	if !ok || user.Role != "admin" {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	animalID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		h.ErrorResponse(w, "Неверный ID животного", http.StatusBadRequest)
		return
	}

	animal, err := h.repo.GetAnimalByID(animalID)
	if err != nil {
		h.ErrorResponse(w, "Ошибка при получении данных животного", http.StatusInternalServerError)
		return
	}

	data := struct {
		User   models.User
		Animal models.Animal
	}{
		User:   user,
		Animal: *animal,
	}

	h.RenderTemplate(w, "edit_animal.html", data)
}

func (h *Handlers) EditAnimal(w http.ResponseWriter, r *http.Request) {
	user, ok := h.GetUserFromSession(w, r)
	if !ok || user.Role != "admin" {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	animalID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		h.ErrorResponse(w, "Неверный ID животного", http.StatusBadRequest)
		return
	}

	if r.Method == http.MethodPost {
		age, err := strconv.Atoi(r.FormValue("age"))
		if err != nil {
			h.ErrorResponse(w, "Неверный возраст", http.StatusBadRequest)
			return
		}

		available := r.FormValue("available") == "on"

		animal := models.Animal{
			ID:          animalID,
			Name:        r.FormValue("name"),
			Type:        r.FormValue("type"),
			Breed:       r.FormValue("breed"),
			Age:         age,
			Description: r.FormValue("description"),
			ImageURL:    r.FormValue("image_url"),
			Available:   available,
		}

		if err := h.repo.UpdateAnimal(&animal); err != nil {
			h.ErrorResponse(w, "Ошибка при обновлении данных животного", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/admin", http.StatusFound)
		return
	}

	animal, err := h.repo.GetAnimalByID(animalID)
	if err != nil {
		h.ErrorResponse(w, "Ошибка при получении данных животного", http.StatusInternalServerError)
		return
	}

	data := struct {
		User   models.User
		Animal models.Animal
	}{
		User:   user,
		Animal: *animal,
	}

	h.RenderTemplate(w, "edit_animal.html", data)
}

func (h *Handlers) CreateAnimal(w http.ResponseWriter, r *http.Request) {
	user, ok := h.GetUserFromSession(w, r)
	if !ok || user.Role != "admin" {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	if r.Method == http.MethodPost {
		age, err := strconv.Atoi(r.FormValue("age"))
		if err != nil {
			h.ErrorResponse(w, "Неверный возраст", http.StatusBadRequest)
			return
		}

		available := r.FormValue("available") == "on"

		animal := models.Animal{
			Name:        r.FormValue("name"),
			Type:        r.FormValue("type"),
			Breed:       r.FormValue("breed"),
			Age:         age,
			Description: r.FormValue("description"),
			ImageURL:    r.FormValue("image_url"),
			Available:   available,
		}

		if err := h.repo.CreateAnimal(&animal); err != nil {
			h.ErrorResponse(w, "Ошибка при создании животного", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/admin", http.StatusFound)
	}
}

func (h *Handlers) DeleteAnimal(w http.ResponseWriter, r *http.Request) {
	_, ok := h.GetUserFromSession(w, r)
	if !ok {
		return
	}

	animalID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		h.ErrorResponse(w, "Неверный ID животного", http.StatusBadRequest)
		return
	}

	if err := h.repo.DeleteAnimal(animalID); err != nil {
		h.ErrorResponse(w, "Ошибка при удалении животного", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin", http.StatusFound)
}
