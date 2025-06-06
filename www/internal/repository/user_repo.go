package repository

import (
	"database/sql"
	"priyutik/internal/models"
)

type UserRepo struct {
	db *sql.DB
}

func (r *UserRepo) CreateUser(user *models.User) error {
	_, err := r.db.Exec(
		"INSERT INTO users (name, phone, email, password, role) VALUES ($1, $2, $3, $4, $5)",
		user.Name, user.Phone, user.Email, user.Password, user.Role,
	)
	return err
}

func (r *UserRepo) GetUserByPhone(phone string) (*models.User, error) {
	var user models.User
	err := r.db.QueryRow(
		"SELECT id, name, phone, email, password, role FROM users WHERE phone = $1",
		phone,
	).Scan(&user.ID, &user.Name, &user.Phone, &user.Email, &user.Password, &user.Role)
	return &user, err
}

func (r *UserRepo) GetUserByID(id int) (*models.User, error) {
	var user models.User
	err := r.db.QueryRow(
		"SELECT id, name, phone, email, role FROM users WHERE id = $1",
		id,
	).Scan(&user.ID, &user.Name, &user.Phone, &user.Email, &user.Role)
	return &user, err
}
