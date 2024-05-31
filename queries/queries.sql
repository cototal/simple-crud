-- SelectOneTask
SELECT * FROM Task
WHERE id = ?

-- SelectAllTasks
SELECT * FROM Task

-- CreateOneTask
INSERT INTO Task (name)
VALUES (?)

-- UpdateOneTask
UPDATE Task
SET name = ?
WHERE id = ?

-- DeleteOneTask
DELETE FROM Task
WHERE id = ?
