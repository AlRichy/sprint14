package main

import (
	"fmt"
	"net/http"
	"time"
)

func NextDateHandler(res http.ResponseWriter, req *http.Request) {
	nowParam := req.URL.Query().Get("now")
	date := req.URL.Query().Get("date")
	repeat := req.URL.Query().Get("repeat")

	now, err := time.Parse(layout, nowParam)
	if err != nil {
		http.Error(res, "Неправильный формат парамeтра now", http.StatusBadRequest)
		return
	}

	newDate, err := NextDate(now, date, repeat)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	res.Header().Set("Content-Type", "text/plain")
	res.WriteHeader(http.StatusOK)
	_, err = res.Write([]byte(newDate))
	if err != nil {
		fmt.Printf("Error in writing a response for /api/nextdate GET request,\n %v", err)
		return
	}
}

func GetTaskHandler(res http.ResponseWriter, req *http.Request) {
	taskID, err := parseAndValidateID(req.URL.Query().Get("id"))
	if err != nil {
		SendErrorResponse(res, err.Error(), http.StatusBadRequest)
		return
	}

	task, err := GetTaskByID(taskID)
	if err != nil {
		SendErrorResponse(res, "задача с указанным id не найдена", http.StatusNotFound)
		return
	}

	sendJSONResponse(res, http.StatusOK, task)
}

func GetTasksHandler(res http.ResponseWriter, req *http.Request) {
	searchTerm := req.URL.Query().Get("search")
	limit := 50

	tasks, err := GetTasks(searchTerm, limit)
	if err != nil {
		SendErrorResponse(res, "ошибка запроса к базе данных", http.StatusInternalServerError)
		return
	}

	if len(tasks) == 0 {
		tasks = []Task{}
	}

	sendJSONResponse(res, http.StatusOK, map[string][]Task{"tasks": tasks})
}

func AddTaskHandler(res http.ResponseWriter, req *http.Request) {
	var task Task
	if err := parseRequestBody(req, &task); err != nil {
		SendErrorResponse(res, "ошибка декодирования JSON", http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		SendErrorResponse(res, "отсутствует обязательное поле title", http.StatusBadRequest)
		return
	}

	if err := validateAndUpdateDate(&task); err != nil {
		SendErrorResponse(res, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := AddTask(task)
	if err != nil {
		SendErrorResponse(res, "ошибка запроса к базе данных", http.StatusInternalServerError)
		return
	}

	sendJSONResponse(res, http.StatusOK, IDResponse{ID: id})
}

func PutTaskHandler(res http.ResponseWriter, req *http.Request) {
	var taskUpdates map[string]interface{}
	if err := parseRequestBody(req, &taskUpdates); err != nil {
		SendErrorResponse(res, "ошибка декодирования JSON", http.StatusBadRequest)
		return
	}

	_, err := validateAndExtractID(taskUpdates)
	if err != nil {
		SendErrorResponse(res, err.Error(), http.StatusBadRequest)
		return
	}

	if err := validateTaskUpdates(taskUpdates); err != nil {
		SendErrorResponse(res, err.Error(), http.StatusBadRequest)
		return
	}

	if _, err := UpdateTask(taskUpdates); err != nil {
		SendErrorResponse(res, "ошибка запроса к базе данных", http.StatusInternalServerError)
		return
	}

	sendJSONResponse(res, http.StatusOK, map[string]interface{}{})
}

func DeleteTaskHandler(res http.ResponseWriter, req *http.Request) {
	taskID, err := parseAndValidateID(req.URL.Query().Get("id"))
	if err != nil {
		SendErrorResponse(res, err.Error(), http.StatusBadRequest)
		return
	}

	if err := deleteTaskIfExists(taskID); err != nil {
		SendErrorResponse(res, err.Error(), http.StatusInternalServerError)
		return
	}

	sendJSONResponse(res, http.StatusOK, map[string]interface{}{})
}

func DoneTaskHandler(res http.ResponseWriter, req *http.Request) {
	taskID, err := parseAndValidateID(req.URL.Query().Get("id"))
	if err != nil {
		SendErrorResponse(res, err.Error(), http.StatusBadRequest)
		return
	}

	task, err := GetTaskByID(taskID)
	if err != nil {
		SendErrorResponse(res, "задача с указанным id не найдена", http.StatusNotFound)
		return
	}

	if task.Repeat == "" {
		if _, err := DeleteTask(taskID); err != nil {
			SendErrorResponse(res, "ошибка запроса к базе данных", http.StatusInternalServerError)
			return
		}
	} else {
		if err := markTaskAsDone(taskID, task); err != nil {
			SendErrorResponse(res, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	sendJSONResponse(res, http.StatusOK, map[string]interface{}{})
}
