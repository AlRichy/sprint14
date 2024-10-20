package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"alrichy/final-project-todo/config"
	"alrichy/final-project-todo/repeat"
	"alrichy/final-project-todo/utils/storage"
	"alrichy/final-project-todo/utils/task"
)

type ResponseJson struct {
	ID    string `json:"id,omitempty"`
	Error string `json:"error,omitempty"`
}

func NextDateHandler(w http.ResponseWriter, r *http.Request) {

	strnow := r.URL.Query().Get("now")
	date := r.URL.Query().Get("date")
	strRepeat := r.URL.Query().Get("repeat")

	now, err := time.Parse(config.Layout, strnow)
	if err != nil {
		log.Fatal(err)
	}
	nextdate, err := repeat.NextDate(now, date, strRepeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	_, err = w.Write([]byte(nextdate))
	if err != nil {
		log.Fatal(err)
	}
}

func PostGetPutTaskHandler(store storage.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var t task.Task
		switch {
		case r.Method == http.MethodPost:
			err := json.NewDecoder(r.Body).Decode(&t)
			if err != nil {
				http.Error(w, `{"error":"ошибка десериализации JSON"}`, http.StatusBadRequest)
				return
			}
			id, err := store.CreateTask(t)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			resp := ResponseJson{ID: id}

			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(resp); err != nil {
				http.Error(w, `{"error":"Ошибка кодирования JSON"}`, http.StatusInternalServerError)
				return
			}

		case r.Method == http.MethodGet:
			id := r.URL.Query().Get("id")
			task, err := store.GetTask(id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(task); err != nil {
				http.Error(w, `{"error":"Ошибка кодирования JSON"}`, http.StatusInternalServerError)
				return
			}

		case r.Method == http.MethodPut:
			err := json.NewDecoder(r.Body).Decode(&t)
			if err != nil {
				http.Error(w, `{"error":"ошибка десериализации JSON"}`, http.StatusBadRequest)
				return
			}
			err = store.UpdateTask(t)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(map[string]string{}); err != nil {
				http.Error(w, `{"error":"Ошибка кодирования JSON"}`, http.StatusInternalServerError)
				return
			}

		case r.Method == http.MethodDelete:
			id := r.URL.Query().Get("id")
			err := store.DeleteTask(id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(map[string]string{}); err != nil {
				http.Error(w, `{"error":"Ошибка кодирования JSON"}`, http.StatusInternalServerError)
				return
			}
		}
	}
}

func GetTasksHandler(store storage.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		search := r.URL.Query().Get("search")
		tasks, err := store.SearchTask(search)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		response := map[string][]task.Task{
			"tasks": tasks,
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, `{"error":"Ошибка кодирования JSON"}`, http.StatusInternalServerError)
			return
		}
	}
}

func DoneTaskHandler(store storage.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		err := store.DoneTask(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]string{}); err != nil {
			http.Error(w, `{"error":"Ошибка кодирования JSON"}`, http.StatusInternalServerError)
			return
		}
	}
}
