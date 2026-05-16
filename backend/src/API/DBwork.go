package API

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool

func DBAllTasks(ctx context.Context) ([]Task, error) {
	rows, err := Pool.Query(ctx, `
    SELECT 
        t.id, 
        t.title, 
        t.description,
        p.priority_name,
        s.status_name,
        t.deadline,
        t.created_at,
        t.updated_at,
        teams.team_name,
        creator.user_name AS created_by,
        assignee.user_name AS assigned_to
    FROM tasks t
    JOIN priorities p ON p.id = t.priority_id
    JOIN statuses s ON s.id = t.status_id
    JOIN teams ON teams.id = t.team_id
    LEFT JOIN users creator ON creator.id = t.created_by
    LEFT JOIN users assignee ON assignee.id = t.assigned_to
    ORDER BY t.id
`)
	if err != nil {
		return nil, fmt.Errorf("query all tasks: %w", err)
	}
	defer rows.Close()
	return scanTasks(rows)
}

func scanTasks(rows pgx.Rows) ([]Task, error) {
	var tasks []Task
	for rows.Next() {
		var t Task
		err := rows.Scan(
			&t.ID,
			&t.Title,
			&t.Description,
			&t.PriorityName,
			&t.StatusName,
			&t.Deadline,
			&t.CreatedAt,
			&t.UpdatedAt,
			&t.TeamName,
			&t.CreatedBy,
			&t.AssignedTo,
		)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, rows.Err()
}
