package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/omeg21/project-repo/internal/config"
	"github.com/omeg21/project-repo/internal/controller"
)

func routes(app *config.AppConfig) http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(Nosurf)
	mux.Use(SessionLoad)

	mux.Get("/", controller.Repo.Home)
	mux.Get("/about", controller.Repo.About)
	mux.Get("/primo", controller.Repo.Primo)
	mux.Get("/jojo", controller.Repo.Jojo)
	mux.Get("/contact", controller.Repo.Contact)
	mux.Get("/reservation", controller.Repo.Reservation)
	mux.Post("/reservation", controller.Repo.PostReservation)
	mux.Get("/Search-availability", controller.Repo.Availability)
	mux.Post("/Search-availability", controller.Repo.PostAvailability)
	mux.Post("/Search-availability-json", controller.Repo.AvailabilityJSON)
	mux.Get("/ReservationSummary",controller.Repo.ReservationSummary)
	mux.Get("/choose-room/{id}",controller.Repo.ChooseRoom)
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))
	return mux
}
