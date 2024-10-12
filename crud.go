package main

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

func AddTask(task Task) (int64, error) {
	result, err := store.db.Exec(
		"INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)",
		task.Date, task.Title, task.Comment, task.Repeat,
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func GetTasks(searchTerm string, limit int) ([]Task, error) {
	query := "SELECT id, date, title, comment, repeat FROM scheduler"
	args := []interface{}{}

	parsedDate, dateErr := time.Parse("02.01.2006", searchTerm)
	switch {
	case dateErr == nil:
		formattedDate := parsedDate.Format("20060102")
		query += " WHERE date = ? ORDER BY date LIMIT ?"
		args = append(args, formattedDate, limit)
	case searchTerm != "":
		query += " WHERE title LIKE ? OR comment LIKE ? ORDER BY date LIMIT ?"
		searchTerm = "%" + searchTerm + "%"
		args = append(args, searchTerm, searchTerm, limit)
	default:
		query += " ORDER BY date LIMIT ?"
		args = append(args, limit)
	}

	rows, err := store.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		err = rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func GetTaskByID(id string) (*Task, error) {
	row := store.db.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?", id)
	var task Task
	err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("task not found")
		}
		return nil, err
	}
	return &task, nil
}

func UpdateTask(taskUpdates map[string]interface{}) (int64, error) {
	query := "UPDATE scheduler SET "
	args := []interface{}{}
	i := 0

	for key, value := range taskUpdates {
		if key != "id" {
			if i > 0 {
				query += ", "
			}
			query += fmt.Sprintf("%s = ?", key)
			args = append(args, value)
			i++
		}
	}

	query += " WHERE id = ?"
	args = append(args, taskUpdates["id"])

	result, err := store.db.Exec(query, args...)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

func DeleteTask(id string) (int64, error) {
	result, err := store.db.Exec("DELETE FROM scheduler WHERE id = ?", id)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

func MarkTaskAsDone(id, date string) error {
	_, err := store.db.Exec("UPDATE scheduler SET date = ? WHERE id = ?", date, id)
	return err
}
