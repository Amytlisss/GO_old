package app

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"priyutik/internal/config"
	"priyutik/internal/handlers"
	"priyutik/internal/repository"

	"github.com/gorilla/sessions"
	_ "github.com/lib/pq"
)

type Application struct {
	DB       *sql.DB
	Store    *sessions.CookieStore
	Config   *config.Config
	Handlers *handlers.Handlers
}

var db *sql.DB

func New(cfg *config.Config) (*Application, error) {
	db, err := sql.Open("postgres", cfg.DB.URI)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Инициализация репозитория и создание таблиц
	repo := repository.NewRepository(db)
	if err := repo.InitDB(); err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	store := sessions.NewCookieStore([]byte(cfg.Session.SecretKey))
	handlers := handlers.NewHandlers(repo, store, cfg)

	return &Application{
		DB:       db,
		Store:    store,
		Config:   cfg,
		Handlers: handlers,
	}, nil
}

// Остальной код остается без изменений
func (a *Application) Run() error {
	a.Handlers.RegisterRoutes()
	log.Printf("Starting server on :%s", a.Config.Server.Port)
	return http.ListenAndServe(":"+a.Config.Server.Port, nil)
}

func (a *Application) Close() {
	if err := a.DB.Close(); err != nil {
		log.Printf("Error closing database: %v", err)
	}
}
