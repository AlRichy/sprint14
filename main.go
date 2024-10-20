package main

import (
	"fmt"
	"net/http"
	"os"

	"alrichy/final-project-todo/auth"
	"alrichy/final-project-todo/config"
	"alrichy/final-project-todo/db"
	"alrichy/final-project-todo/utils/handler"
	"alrichy/final-project-todo/utils/storage"

	"github.com/go-chi/chi"
)

var Password string

func main() {
	env := config.GetEnv()
	Password = os.Getenv("TODO_PASSWORD")
	fmt.Println("Port:", env.Port)

	dataBase := db.EnsureDB()
	defer dataBase.Close()
	fmt.Println("DB initialized")

	store := storage.NewStore(dataBase)
	fmt.Println("Store initialized")

	r := RegisterRoutes(store)
	fmt.Println("Routes initialized")

	fmt.Println("Server start...")
	err := http.ListenAndServe(":"+env.Port, r)
	if err != nil {
		fmt.Println("Error on starting server:\n", err)
	}
}

func RegisterRoutes(store storage.Store) *chi.Mux {
	r := chi.NewRouter()
	r.Handle("/*", http.FileServer(http.Dir("./web")))
	r.Get("/api/nextdate", handler.NextDateHandler)
	r.Post("/api/task", auth.AuthMiddleware(handler.PostGetPutTaskHandler(store)))
	r.Get("/api/tasks", auth.AuthMiddleware(handler.GetTasksHandler(store)))
	r.Get("/api/task", handler.PostGetPutTaskHandler(store))
	r.Put("/api/task", handler.PostGetPutTaskHandler(store))
	r.Post("/api/task/done", auth.AuthMiddleware(handler.DoneTaskHandler(store)))
	r.Delete("/api/task", handler.PostGetPutTaskHandler(store))
	r.Post("/api/signin", auth.SigninHandler)
	return r
}
