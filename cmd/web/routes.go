package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *Config) routes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(app.SessionLoad)

	mux.Get("/", app.HomePage)
	mux.Get("/login", app.LoginPage)
	mux.Post("/login", app.PostLoginPage)
	mux.Get("/logout", app.Logout)
	mux.Get("/register", app.RegisterPage)
	mux.Post("/register", app.PostRegister)
	mux.Get("/activate-account", app.ActivateAccount)
	mux.Get("/test-mail", func(w http.ResponseWriter, r *http.Request) {
		m := Mail{
			Domain:      "localhost",
			Host:        "local",
			Port:        1025,
			Encryption:  "none",
			FromAddress: "test@test.com",
			FromName:    "test",
			ErrorChan:   make(chan error),
		}
		mail := Message{
			To:      "me@test.com",
			Subject: "test",
			Data:    "test hello",
		}
		m.sendMail(mail, make(chan error))
	})

	return mux
}
