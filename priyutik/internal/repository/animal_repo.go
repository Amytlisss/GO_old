package repository

import (
	"database/sql"
	"priyutik/internal/models"
)

type AnimalRepo struct {
	db *sql.DB
}

func (r *AnimalRepo) GetAllAnimals() ([]models.Animal, error) {
	rows, err := r.db.Query(
		"SELECT id, name, type, breed, age, description, image_url, available FROM animals",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var animals []models.Animal
	for rows.Next() {
		var a models.Animal
		if err := rows.Scan(&a.ID, &a.Name, &a.Type, &a.Breed, &a.Age, &a.Description, &a.ImageURL, &a.Available); err != nil {
			return nil, err
		}
		animals = append(animals, a)
	}
	return animals, nil
}

func (r *AnimalRepo) GetAnimalByID(id int) (*models.Animal, error) {
	var a models.Animal
	err := r.db.QueryRow(
		"SELECT id, name, type, breed, age, description, image_url, available FROM animals WHERE id = $1",
		id,
	).Scan(&a.ID, &a.Name, &a.Type, &a.Breed, &a.Age, &a.Description, &a.ImageURL, &a.Available)
	return &a, err
}

func (r *AnimalRepo) CreateAnimal(animal *models.Animal) error {
	_, err := r.db.Exec(
		"INSERT INTO animals (name, type, breed, age, description, image_url, available) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		animal.Name, animal.Type, animal.Breed, animal.Age, animal.Description, animal.ImageURL, animal.Available,
	)
	return err
}

func (r *AnimalRepo) UpdateAnimal(animal *models.Animal) error {
	_, err := r.db.Exec(
		"UPDATE animals SET name = $1, type = $2, breed = $3, age = $4, description = $5, image_url = $6, available = $7 WHERE id = $8",
		animal.Name, animal.Type, animal.Breed, animal.Age, animal.Description, animal.ImageURL, animal.Available, animal.ID,
	)
	return err
}

func (r *AnimalRepo) DeleteAnimal(id int) error {
	_, err := r.db.Exec(
		"DELETE FROM animals WHERE id = $1",
		id,
	)
	return err
}
