package app

import (
	"database/sql"
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
	var err error
	connStr := "user=postgres password=0000 dbname=priyutik sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	store := sessions.NewCookieStore([]byte(cfg.Session.SecretKey))

	repo := repository.NewRepository(db)
	handlers := handlers.NewHandlers(repo, store, cfg)

	return &Application{
		DB:       db,
		Store:    store,
		Config:   cfg,
		Handlers: handlers,
	}, nil
}

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
