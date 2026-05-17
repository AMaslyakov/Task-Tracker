package API

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool

func DBAllTasks(ctx context.Context) ([]Task, error) {
	rows, err := Pool.Query(ctx, `
    SELECT 
        t.id, 
        t.title, 
        COALESCE(t.description, ''),
        p.priority_name,
        s.status_name,
        t.deadline,
        t.created_at,
        t.updated_at,
        t.team_id,
        teams.team_name,
        creator.user_name AS created_by,
        COALESCE(assignee.user_name, '') AS assigned_to
    FROM tasks t
    JOIN priorities p ON p.id = t.priority_id
    JOIN statuses s ON s.id = t.status_id
    JOIN teams ON teams.id = t.team_id
    LEFT JOIN users creator ON creator.id = t.created_by
    LEFT JOIN users assignee ON assignee.id = t.assigned_to
    ORDER BY t.id`)
	if err != nil {
		return nil, fmt.Errorf("query all tasks: %w", err)
	}
	defer rows.Close()
	return scanTasks(rows)
}

func DBTasksByTeam(ctx context.Context, teamID int) ([]Task, error) {
	rows, err := Pool.Query(ctx, `
    SELECT 
        t.id, 
        t.title, 
        COALESCE(t.description, ''),
        p.priority_name,
        s.status_name,
        t.deadline,
        t.created_at,
        t.updated_at,
        t.team_id,
        teams.team_name,
        creator.user_name AS created_by,
        COALESCE(assignee.user_name, '') AS assigned_to
    FROM tasks t
    JOIN priorities p ON p.id = t.priority_id
    JOIN statuses s ON s.id = t.status_id
    JOIN teams ON teams.id = t.team_id
    LEFT JOIN users creator ON creator.id = t.created_by
    LEFT JOIN users assignee ON assignee.id = t.assigned_to
    WHERE t.team_id = $1
    ORDER BY t.id`, teamID)
	if err != nil {
		return nil, fmt.Errorf("query tasks by team: %w", err)
	}
	defer rows.Close()
	return scanTasks(rows)
}

func DBTaskByID(ctx context.Context, taskID int) (Task, error) {
	rows, err := Pool.Query(ctx, `
    SELECT 
        t.id, 
        t.title, 
        COALESCE(t.description, ''),
        p.priority_name,
        s.status_name,
        t.deadline,
        t.created_at,
        t.updated_at,
        t.team_id,
        teams.team_name,
        creator.user_name AS created_by,
        COALESCE(assignee.user_name, '') AS assigned_to
    FROM tasks t
    JOIN priorities p ON p.id = t.priority_id
    JOIN statuses s ON s.id = t.status_id
    JOIN teams ON teams.id = t.team_id
    LEFT JOIN users creator ON creator.id = t.created_by
    LEFT JOIN users assignee ON assignee.id = t.assigned_to
    WHERE t.id = $1`, taskID)
	if err != nil {
		return Task{}, fmt.Errorf("query task by id: %w", err)
	}
	defer rows.Close()

	tasks, err := scanTasks(rows)
	if err != nil {
		return Task{}, err
	}
	if len(tasks) == 0 {
		return Task{}, pgx.ErrNoRows
	}
	return tasks[0], nil
}

func DBAllTeams(ctx context.Context) ([]Team, error) {
	query := `
        SELECT 
            t.id,
            t.team_name,
            COALESCE(t.description, ''),
            t.config_dashboard,
            t.created_at::TEXT,
            t.updated_at::TEXT,
            COALESCE((SELECT ARRAY_AGG(title) FROM tasks WHERE team_id = t.id), '{}') AS tasks,
            COALESCE((SELECT ARRAY_AGG(u.user_name) FROM team_members tm JOIN users u ON u.id = tm.user_id WHERE tm.team_id = t.id AND tm.is_active = TRUE), '{}') AS members
        FROM teams t
        ORDER BY t.id
    `
	rows, err := Pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query teams: %w", err)
	}
	defer rows.Close()

	var teams []Team
	for rows.Next() {
		var team Team
		var configDashboard []byte
		err := rows.Scan(
			&team.ID,
			&team.Name,
			&team.Description,
			&configDashboard,
			&team.CreatedAt,
			&team.UpdatedAt,
			&team.Tasks,
			&team.Members,
		)
		if err != nil {
			return nil, fmt.Errorf("scan team: %w", err)
		}
		if err := json.Unmarshal(configDashboard, &team.ConfigDashboard); err != nil {
			return nil, fmt.Errorf("parse team config dashboard: %w", err)
		}
		if team.Tasks == nil {
			team.Tasks = []string{}
		}
		if team.Members == nil {
			team.Members = []string{}
		}
		if team.ConfigDashboard == nil {
			team.ConfigDashboard = []string{}
		}
		teams = append(teams, team)
	}
	return teams, rows.Err()
}

func DBTeamByID(ctx context.Context, teamID int) (Team, error) {
	query := `
        SELECT 
            t.id,
            t.team_name,
            COALESCE(t.description, ''),
            t.config_dashboard,
            t.created_at::TEXT,
            t.updated_at::TEXT,
            COALESCE((SELECT ARRAY_AGG(title) FROM tasks WHERE team_id = t.id), '{}') AS tasks,
            COALESCE((SELECT ARRAY_AGG(u.user_name) FROM team_members tm JOIN users u ON u.id = tm.user_id WHERE tm.team_id = t.id AND tm.is_active = TRUE), '{}') AS members
        FROM teams t
        WHERE t.id = $1
    `
	var team Team
	var configDashboard []byte
	err := Pool.QueryRow(ctx, query, teamID).Scan(
		&team.ID,
		&team.Name,
		&team.Description,
		&configDashboard,
		&team.CreatedAt,
		&team.UpdatedAt,
		&team.Tasks,
		&team.Members,
	)
	if err != nil {
		return Team{}, fmt.Errorf("query team by id: %w", err)
	}
	if err := json.Unmarshal(configDashboard, &team.ConfigDashboard); err != nil {
		return Team{}, fmt.Errorf("parse team config dashboard: %w", err)
	}
	if team.Tasks == nil {
		team.Tasks = []string{}
	}
	if team.Members == nil {
		team.Members = []string{}
	}
	if team.ConfigDashboard == nil {
		team.ConfigDashboard = []string{}
	}
	return team, nil
}

func DBInsertTask(ctx context.Context) error {
	return nil
}

func CorrectSession(ctx context.Context, SessionId string) (int, bool) {
	var userId int
	err := Pool.QueryRow(ctx, "SELECT UserId FROM Sessions WHERE SessionId = $1", SessionId).Scan(&userId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, false
		}
		log.Printf("DB error in CorrectSession: %v", err)
		return 0, false
	}
	return userId, true
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
			&t.TeamID,
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
