package handlers

import (
	"net/http"
	"strconv"

	"priyutik/internal/models"
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

func (h *Handlers) AddAnimal(w http.ResponseWriter, r *http.Request) {
	user, _ := h.GetUserFromSession(w, r)

	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			h.ErrorResponse(w, "Ошибка при парсинге формы", http.StatusBadRequest)
			return
		}

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
			h.ErrorResponse(w, "Ошибка при добавлении животного", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		return
	}

	h.RenderTemplate(w, "add_animal.html", struct {
		User models.User
	}{
		User: user,
	})
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
