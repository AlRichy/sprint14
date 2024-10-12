package main

import (
	"fmt"
	"net/http"
	"os"
)

func RunApp() {
	// Task 1
	runOnPort()
}

func runOnPort() {
	// 1*. Если существует переменная окружения TODO_PORT, сервер при старте будет слушать порт со значением этой переменной

	todoPort, exists := os.LookupEnv("TODO_PORT")
	if !exists {
		todoPort = "7540"
	}

	router := RegisterRoutes()

	// 1. Файловый сервер будет возвращать файлы из поддиректории web
	fs := http.FileServer(http.Dir("./web"))
	router.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	})
	http.ListenAndServe(fmt.Sprintf(":%s", todoPort), router)
}
