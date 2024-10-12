package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func parseRequestBody(req *http.Request, target interface{}) error {
	var buf bytes.Buffer
	if _, err := buf.ReadFrom(req.Body); err != nil {
		return fmt.Errorf("ошибка чтения тела запроса")
	}
	return json.Unmarshal(buf.Bytes(), target)
}

func validateAndUpdateDate(task *Task) error {
	var dateInTime time.Time
	var err error

	if task.Date != "" {
		dateInTime, err = time.Parse(layout, task.Date)
		if err != nil {
			return fmt.Errorf("недопустимый формат date")
		}
	} else {
		task.Date = time.Now().Format(layout)
		dateInTime = time.Now()
	}

	if time.Now().After(dateInTime) {
		if task.Repeat != "" {
			task.Date, err = NextDate(time.Now(), task.Date, task.Repeat)
			if err != nil {
				return err
			}
		} else {
			task.Date = time.Now().Format(layout)
		}
	}

	return nil
}

func parseAndValidateID(idStr string) (string, error) {
	if idStr == "" {
		return "", fmt.Errorf("не передан идентификатор")
	}

	if _, err := strconv.ParseInt(idStr, 10, 64); err != nil {
		return "", fmt.Errorf("id должен быть числом")
	}
	return idStr, nil
}

func validateAndExtractID(taskUpdates map[string]interface{}) (string, error) {
	id, ok := taskUpdates["id"].(string)
	if !ok || id == "" {
		return "", fmt.Errorf("отсутствует обязательное поле id")
	}

	if _, err := strconv.ParseInt(id, 10, 64); err != nil {
		return "", fmt.Errorf("id должен быть числом")
	}

	if _, err := GetTaskByID(id); err != nil {
		return "", fmt.Errorf("задача с указанным id не найдена")
	}

	return id, nil
}

func validateTaskUpdates(taskUpdates map[string]interface{}) error {
	date, dateOk := taskUpdates["date"].(string)
	if title, ok := taskUpdates["title"].(string); !ok || strings.TrimSpace(title) == "" {
		return fmt.Errorf("отсутствует обязательное поле title")
	}

	if !dateOk || strings.TrimSpace(date) == "" {
		return fmt.Errorf("отсутствует обязательное поле date")
	}

	if _, err := time.Parse(layout, date); err != nil {
		return fmt.Errorf("недопустимый формат date")
	}

	if repeat, ok := taskUpdates["repeat"].(string); ok && strings.TrimSpace(repeat) != "" {
		repeatParts := strings.SplitN(repeat, " ", 2)
		repeatType := repeatParts[0]
		if !isValidRepeatType(repeatType) {
			return fmt.Errorf("недопустимый символ")
		}
	}

	return nil
}

func isValidRepeatType(repeatType string) bool {
	validTypes := []string{"d", "w", "m", "y"}
	for _, v := range validTypes {
		if repeatType == v {
			return true
		}
	}
	return false
}

func deleteTaskIfExists(taskID string) error {
	if _, err := GetTaskByID(taskID); err != nil {
		return fmt.Errorf("задача с указанным id не найдена")
	}

	if _, err := DeleteTask(taskID); err != nil {
		return fmt.Errorf("ошибка запроса к базе данных")
	}

	return nil
}

func markTaskAsDone(taskID string, task *Task) error {
	parsedDate, err := time.Parse(layout, task.Date)
	if err != nil {
		return fmt.Errorf("недопустимый формат date")
	}

	if parsedDate.Format(layout) == time.Now().Format(layout) {
		parsedDate = parsedDate.AddDate(0, 0, -1)
		task.Date, err = NextDate(parsedDate, task.Date, task.Repeat)
	} else {
		task.Date, err = NextDate(time.Now(), task.Date, task.Repeat)
	}

	if err != nil {
		return err
	}

	if err := MarkTaskAsDone(taskID, task.Date); err != nil {
		return fmt.Errorf("ошибка при обновлении задачи")
	}

	return nil
}
