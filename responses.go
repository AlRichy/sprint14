package main

import (
	"encoding/json"
	"net/http"
)

func SendErrorResponse(res http.ResponseWriter, errorMessage string, statusCode int) {
	resp := ErrorResponse{Error: errorMessage}
	respBytes, err := json.Marshal(resp)
	if err != nil {
		http.Error(res, "ошибка при сериализации ответа", http.StatusInternalServerError)
		return
	}
	http.Error(res, string(respBytes), statusCode)
}

func sendJSONResponse(res http.ResponseWriter, statusCode int, data interface{}) {
	respBytes, err := json.Marshal(data)
	if err != nil {
		SendErrorResponse(res, "ошибка при сериализации ответа", http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(statusCode)

	_, err = res.Write(respBytes)
	if err != nil {
		SendErrorResponse(res, "ошибка записи ответа", http.StatusInternalServerError)
		return
	}
}
