package main

import (
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes() *chi.Mux {
	r := chi.NewRouter()

	r.Get("/api/nextdate", NextDateHandler)
	r.Post("/api/task", Auth(AddTaskHandler))
	r.Get("/api/tasks", Auth(GetTasksHandler))
	r.Get("/api/task", Auth(GetTaskHandler))
	r.Put("/api/task", Auth(PutTaskHandler))
	r.Delete("/api/task", Auth(DeleteTaskHandler))
	r.Post("/api/task/done", Auth(DoneTaskHandler))
	r.Post("/api/signin", SignInHandler)
	return r
}
