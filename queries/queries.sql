-- SelectOneTask
SELECT * FROM task
WHERE id = ?

-- SelectAllTasks
SELECT * FROM task

-- CreateOneTask
INSERT INTO task (name)
VALUES (?)

-- UpdateOneTask
UPDATE task
SET name = ?
WHERE id = ?

-- DeleteOneTask
DELETE FROM task
WHERE id = ?
