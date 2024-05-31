package queries

func SelectOneTask() string {
	return "SELECT * FROM task WHERE id = ?"
}

func SelectAllTasks() string {
	return "SELECT * FROM task"
}

func CreateOneTask() string {
	return "INSERT INTO task (name) VALUES (?)"
}

func UpdateOneTask() string {
	return "UPDATE task SET name = ? WHERE id = ?"
}

func DeleteOneTask() string {
	return "DELETE FROM task WHERE id = ?"
}