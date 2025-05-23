package repository

import (
	"database/sql"
	"log"
)

type Repository struct {
	db *sql.DB
	UserRepo
	MeetingRepo
	AnimalRepo
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db:          db,
		UserRepo:    UserRepo{db: db},
		MeetingRepo: MeetingRepo{db: db},
		AnimalRepo:  AnimalRepo{db: db},
	}
}

func (r *Repository) InitDB() error {
	_, err := r.db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			name TEXT,
			phone TEXT UNIQUE,
			email TEXT UNIQUE,
			password TEXT,
			role TEXT DEFAULT 'user'
		);
		
		CREATE TABLE IF NOT EXISTS meetings (
			id SERIAL PRIMARY KEY,
			user_id INT REFERENCES users(id),
			date TIMESTAMP,
			cancelled BOOLEAN DEFAULT FALSE,
			created_at TIMESTAMP DEFAULT NOW()
		);
		
		CREATE TABLE IF NOT EXISTS animals (
			id SERIAL PRIMARY KEY,
			name TEXT,
			type TEXT,
			breed TEXT,
			age INT,
			description TEXT,
			image_url TEXT,
			available BOOLEAN DEFAULT TRUE
		);
	`)

	if err != nil {
		return err
	}

	// Insert sample animals if table is empty
	var count int
	err = r.db.QueryRow("SELECT COUNT(*) FROM animals").Scan(&count)
	if err == nil && count == 0 {
		_, err = r.db.Exec(`
			INSERT INTO animals (name, type, breed, age, description, image_url) 
			VALUES 
				('Барбос', 'dog', 'Метис', 2, 'Дружелюбный пёс', 'https://example.com/dog1.jpg'),
				('Мышка', 'dog', 'Лайка', 1, 'Активная собака', 'https://example.com/dog2.jpg'),
				('Пушок', 'cat', 'Длинношёрстная', 1, 'Ласковый кот', 'https://example.com/cat1.jpg')
		`)
		if err != nil {
			log.Printf("Error inserting sample animals: %v", err)
		}
	}

	return nil
}
