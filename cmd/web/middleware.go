package main

import (
	"fmt"
	"net/http"

	"github.com/justinas/nosurf"
)

// WriteToConsole does funny
func WriteToConsole(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Landed to the page safely")
		next.ServeHTTP(w, r)
	})
}

// noSurf adds CSRF protection to all POST request
func NoSurf(next http.Handler) http.Handler {
	CSRFHandler := nosurf.New(next)

	CSRFHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   app.InProduction,
		SameSite: http.SameSiteLaxMode,
	})

	return CSRFHandler
}

// SessionLoad Loads and save the sessions on every reques
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}
