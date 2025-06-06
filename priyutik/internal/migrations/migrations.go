package migrations

import (
	"database/sql"
	"fmt"
	"log"
)

type Migration struct {
	Name string
	Up   func(*sql.DB) error
	Down func(*sql.DB) error
}

var migrations = []Migration{
	{
		Name: "create_initial_tables",
		Up: func(db *sql.DB) error {
			_, err := db.Exec(`
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
			return err
		},
		Down: func(db *sql.DB) error {
			_, err := db.Exec(`
				DROP TABLE IF EXISTS animals;
				DROP TABLE IF EXISTS meetings;
				DROP TABLE IF EXISTS users;
			`)
			return err
		},
	},
	{
		Name: "add_sample_animals",
		Up: func(db *sql.DB) error {
			var count int
			err := db.QueryRow("SELECT COUNT(*) FROM animals").Scan(&count)
			if err != nil {
				return fmt.Errorf("error counting animals: %w", err)
			}

			if count == 0 {
				_, err = db.Exec(`
					INSERT INTO animals (name, type, breed, age, description, image_url) 
					VALUES 
						('Барбос', 'dog', 'Метис', 2, 'Дружелюбный пёс', 'https://example.com/dog1.jpg'),
						('Мышка', 'dog', 'Лайка', 1, 'Активная собака', 'https://example.com/dog2.jpg'),
						('Пушок', 'cat', 'Длинношёрстная', 1, 'Ласковый кот', 'https://example.com/cat1.jpg')
				`)
				return err
			}
			return nil
		},
		Down: func(db *sql.DB) error {
			_, err := db.Exec("DELETE FROM animals")
			return err
		},
	},
}

func RunMigrations(db *sql.DB) error {
	// Создаем таблицу для отслеживания выполненных миграций
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS migrations (
			id SERIAL PRIMARY KEY,
			name TEXT UNIQUE,
			applied_at TIMESTAMP DEFAULT NOW()
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Получаем список уже выполненных миграций
	rows, err := db.Query("SELECT name FROM migrations")
	if err != nil {
		return fmt.Errorf("failed to query applied migrations: %w", err)
	}
	defer rows.Close()

	applied := make(map[string]bool)
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return fmt.Errorf("failed to scan migration name: %w", err)
		}
		applied[name] = true
	}

	// Применяем новые миграции
	for _, migration := range migrations {
		if !applied[migration.Name] {
			log.Printf("Applying migration: %s", migration.Name)
			if err := migration.Up(db); err != nil {
				return fmt.Errorf("failed to apply migration %s: %w", migration.Name, err)
			}

			_, err = db.Exec("INSERT INTO migrations (name) VALUES ($1)", migration.Name)
			if err != nil {
				return fmt.Errorf("failed to record migration %s: %w", migration.Name, err)
			}
		}
	}

	return nil
}
