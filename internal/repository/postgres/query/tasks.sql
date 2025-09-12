-- name: GetTaskByID :one
SELECT * FROM tasks WHERE id = $1;

-- name: GetAllTasks :many
SELECT * FROM tasks ORDER BY created_at DESC;

-- name: SaveTask :exec
INSERT INTO tasks (id, title, status, created_at, due_date, priority)
VALUES ($1, $2, $3, $4, $5, $6)
ON CONFLICT (id) DO UPDATE
SET title = EXCLUDED.title,
    status = EXCLUDED.status,
    created_at = EXCLUDED.created_at,
    due_date = EXCLUDED.due_date,
    priority = EXCLUDED.priority;

-- name: DeleteTask :exec
DELETE FROM tasks WHERE id = $1;

-- name: GetTasksByStatus :many
SELECT * FROM tasks WHERE status = $1 ORDER BY created_at DESC;

-- name: GetSortedTasks :many
SELECT * FROM tasks
ORDER BY
    CASE WHEN $1 = 'title' AND $2 = 'asc' THEN title END ASC,
    CASE WHEN $1 = 'title' AND $2 = 'desc' THEN title END DESC,
    CASE WHEN $1 = 'created_at' AND $2 = 'asc' THEN created_at END ASC,
    CASE WHEN $1 = 'created_at' AND $2 = 'desc' THEN created_at END DESC,
    CASE WHEN $1 = 'priority' AND $2 = 'asc' THEN
        CASE priority
            WHEN 'low' THEN 1
            WHEN 'medium' THEN 2
            WHEN 'high' THEN 3
        END
    END ASC,
    CASE WHEN $1 = 'priority' AND $2 = 'desc' THEN
        CASE priority
            WHEN 'low' THEN 1
            WHEN 'medium' THEN 2
            WHEN 'high' THEN 3
        END
    END DESC;

-- name: GetTasksDueToday :many
SELECT * FROM tasks 
WHERE due_date::date = CURRENT_DATE;

-- name: GetTasksDueThisWeek :many
SELECT * FROM tasks 
WHERE due_date >= CURRENT_DATE
  AND due_date < CURRENT_DATE + INTERVAL '7 days';

-- name: GetOverdueTasks :many
SELECT * FROM tasks 
WHERE due_date < CURRENT_DATE
  AND status = 'active';