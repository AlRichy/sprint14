package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v4"
)

var TODO_PASS = "1234"

func SignInHandler(res http.ResponseWriter, req *http.Request) {
	pass := TODO_PASS
	if value, exists := os.LookupEnv("TODO_PASSWORD"); exists {
		pass = value
		TODO_PASS = pass
	}
	if len(pass) == 0 {
		SendErrorResponse(res, "пароль не установлен", http.StatusInternalServerError)
		return
	}

	var buf bytes.Buffer
	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		SendErrorResponse(res, "ошибка чтения тела запроса", http.StatusBadRequest)
		return
	}

	var body map[string]string
	err = json.Unmarshal(buf.Bytes(), &body)
	if err != nil {
		SendErrorResponse(res, "ошибка декодирования JSON", http.StatusBadRequest)
		return
	}

	if body["password"] != pass {
		SendErrorResponse(res, "неверный пароль", http.StatusUnauthorized)
		return
	}

	token := jwt.New(jwt.SigningMethodHS256)
	tokenString, err := token.SignedString([]byte(pass))
	if err != nil {
		SendErrorResponse(res, "ошибка создания токена", http.StatusInternalServerError)
		return
	}

	var resp AuthResponse
	resp.Token = tokenString
	respBytes, err := json.Marshal(resp)
	if err != nil {
		http.Error(res, "ошибка при сериализации ответа", http.StatusInternalServerError)
		return
	}
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	_, err = res.Write(respBytes)
	if err != nil {
		fmt.Println("ошибка записи ответа в HandleSignIn", err)
	}
}
