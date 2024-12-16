package taskrepo

import (
	"context"
	"fmt"
	"kn-assignment/internal/constant"
	"kn-assignment/internal/core/domain"
	errors "kn-assignment/internal/core/error"

	"github.com/georgysavva/scany/v2/pgxscan"
)

func (r *repository) CreateTask(ctx context.Context, task domain.CreateTaskRequest, userId string) error {
	query := `INSERT INTO tasks (title, description, due_date, created_at, created_by, updated_at, updated_by) VALUES ($1, $2, $3, NOW(), $4, NOW(), $4)`
	_, err := r.dbPool.Exec(ctx, query, task.Title, task.Description, task.DueDate, userId)
	if err != nil {
		return errors.NewCustomError(constant.ErrCodeInternalServer)
	}
	return nil
}

func (r *repository) AssignTask(ctx context.Context, taskID, assigneeID string) error {
	query := `UPDATE tasks SET assignee_id = $1 WHERE id = $2`
	_, err := r.dbPool.Exec(ctx, query, assigneeID, taskID)
	if err != nil {
		return errors.NewCustomError(constant.ErrCodeInternalServer)
	}
	return nil
}

func (r *repository) GetTasksByAssignee(ctx context.Context, assigneeID string) ([]domain.Task, error) {
	query := `SELECT * FROM tasks WHERE assignee_id = $1`
	var tasks []domain.Task
	err := pgxscan.Select(ctx, r.dbPool, &tasks, query, assigneeID)
	if err != nil {
		return nil, errors.NewCustomError(constant.ErrCodeInternalServer)
	}
	return tasks, nil
}

func (r *repository) UpdateTaskStatus(ctx context.Context, taskID string, status domain.TaskStatus, userId string) error {
	query := `UPDATE tasks SET status = $1, updated_by = $2, updated_at = NOW() WHERE id = $3`
	_, err := r.dbPool.Exec(ctx, query, status, userId, taskID)
	if err != nil {
		return errors.NewCustomError(constant.ErrCodeInternalServer)
	}
	return nil
}

func (r *repository) GetAllTasks(ctx context.Context, filter map[string]string, sort, order string) ([]domain.Task, error) {
	query := `SELECT * FROM tasks WHERE 1=1`
	args := []interface{}{}
	i := 1
	for k, v := range filter {
		query += fmt.Sprintf(" AND %s = $%d", k, i)
		args = append(args, v)
		i++
	}
	if sort != "" {
		query += fmt.Sprintf(" ORDER BY %s %s", sort, order)
	} else {
		query += " ORDER BY created_at DESC, status ASC"
	}
	var tasks []domain.Task
	err := pgxscan.Select(ctx, r.dbPool, &tasks, query, args...)
	if err != nil {
		return nil, errors.NewCustomError(constant.ErrCodeInternalServer)
	}
	return tasks, nil
}

func (r *repository) GetTaskSummary(ctx context.Context) ([]domain.TaskSummary, error) {
	query := `SELECT assignee_id as employee_id, COUNT(*) as total_tasks, SUM(CASE WHEN status = 'Completed' THEN 1 ELSE 0 END) as completed_tasks FROM tasks WHERE assignee_id IS NOT NULL GROUP BY assignee_id`
	var summaries []domain.TaskSummary
	err := pgxscan.Select(ctx, r.dbPool, &summaries, query)
	if err != nil {
		return nil, errors.NewCustomError(constant.ErrCodeInternalServer)
	}
	return summaries, nil
}

func (r *repository) GetTaskByID(ctx context.Context, taskID string) (domain.Task, error) {
	query := `SELECT * FROM tasks WHERE id = $1`
	var task domain.Task
	err := pgxscan.Get(ctx, r.dbPool, &task, query, taskID)
	if err != nil {
		return domain.Task{}, errors.NewCustomError(constant.ErrCodeInternalServer)
	}
	return task, nil
}

func (r *repository) UpdateTask(ctx context.Context, taskID string, name, description *string) error {
	query := `UPDATE tasks SET`
	args := []interface{}{}
	i := 1

	if name != nil {
		query += fmt.Sprintf(" name = $%d,", i)
		args = append(args, *name)
		i++
	}

	if description != nil {
		query += fmt.Sprintf(" description = $%d,", i)
		args = append(args, *description)
		i++
	}

	query = query[:len(query)-1] // Remove the trailing comma
	query += fmt.Sprintf(" WHERE id = $%d", i)
	args = append(args, taskID)

	_, err := r.dbPool.Exec(ctx, query, args...)
	if err != nil {
		return errors.NewCustomError(constant.ErrCodeInternalServer)
	}

	return nil
}

func (r *repository) DeleteTask(ctx context.Context, taskID string) error {
	query := `DELETE FROM tasks WHERE id = $1`
	_, err := r.dbPool.Exec(ctx, query, taskID)
	if err != nil {
		return errors.NewCustomError(constant.ErrCodeInternalServer)
	}
	return nil
}
