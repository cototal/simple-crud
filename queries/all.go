package queries

func SelectOneTask() string {
	return "SELECT * FROM Task WHERE id = ?"
}

func SelectAllTasks() string {
	return "SELECT * FROM Task"
}

func CreateOneTask() string {
	return "INSERT INTO Task (name) VALUES (?)"
}

func UpdateOneTask() string {
	return "UPDATE Task SET name = ? WHERE id = ?"
}

func DeleteOneTask() string {
	return "DELETE FROM Task WHERE id = ?"
}