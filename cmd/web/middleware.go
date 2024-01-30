package main

import "net/http"

func (app *Configuration) SessionLoad(next http.Handler) http.Handler {
	return app.Session.LoadAndSave(next)
}