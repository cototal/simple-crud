package repos

import (
	"cototal/simple-crud/queries"
	"database/sql"
)

type Task struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func GetAllTasks(db *sql.DB) ([]Task, error) {
	rows, err := db.Query(queries.SelectAllTasks())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := make([]Task, 0, 20)
	for rows.Next() {
		var task Task
		if err := rows.Scan(&task.ID, &task.Name); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func GetTask(db *sql.DB, id int) (Task, error) {
	var task Task
	err := db.QueryRow(queries.SelectOneTask(), id).Scan(
		&task.ID, &task.Name)

	if err != nil {
		return Task{}, err
	}

	return task, nil
}

func CreateTask(db *sql.DB, task *Task) error {
	result, err := db.Exec(queries.CreateOneTask(), task.Name)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	task.ID = int(id)
	return nil
}

func UpdateTask(db *sql.DB, id int, task *Task) error {
	_, err := db.Exec(queries.UpdateOneTask(), task.Name, id)
	if err != nil {
		return err
	}
	task.ID = id
	return nil
}

func DeleteTask(db *sql.DB, id int) error {
	_, err := db.Exec(queries.DeleteOneTask(), id)
	if err != nil {
		return err
	}

	return nil
}
