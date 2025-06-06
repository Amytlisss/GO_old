package repository

import (
	"database/sql"
	"fmt"
	"log"
	"priyutik/internal/migrations"
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
	log.Println("Running database migrations...")
	if err := migrations.RunMigrations(r.db); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}
	log.Println("Migrations completed successfully")
	return nil
}
