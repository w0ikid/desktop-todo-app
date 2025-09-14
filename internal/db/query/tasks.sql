-- name: GetTaskByID :one
SELECT * FROM tasks WHERE id = $1;

-- name: GetAllTasks :many
SELECT * FROM tasks ORDER BY created_at DESC;

-- name: SaveTask :exec
INSERT INTO tasks (id, title, status, created_at, due_date, priority)
VALUES ($1, $2, $3, $4, $5, $6)
ON CONFLICT (id) DO UPDATE
SET title      = EXCLUDED.title,
    status     = EXCLUDED.status,
    created_at = EXCLUDED.created_at,
    due_date   = EXCLUDED.due_date,
    priority   = EXCLUDED.priority;

-- name: DeleteTask :exec
DELETE FROM tasks WHERE id = $1;

-- name: GetTasksByStatus :many
SELECT * FROM tasks
WHERE status = $1
ORDER BY created_at DESC;

-- name: GetTasksDueBetween :many
SELECT * FROM tasks
WHERE due_date >= $1
  AND due_date < $2
ORDER BY due_date ASC;
