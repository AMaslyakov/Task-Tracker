package API

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

var Pool *pgxpool.Pool

var ErrRelatedEntityNotFound = errors.New("related entity not found")

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

func DBInsertTask(ctx context.Context, req CreateTaskRequest) (Task, error) {
	var taskID int
	err := Pool.QueryRow(ctx, `
        WITH refs AS (
            SELECT
                s.id AS status_id,
                p.id AS priority_id,
                tm.id AS team_id,
                creator.id AS created_by,
                assignee.id AS assigned_to
            FROM statuses s
            JOIN priorities p ON p.id = $4
            JOIN teams tm ON tm.id = $6
            JOIN users creator ON creator.id = $7
            LEFT JOIN users assignee ON assignee.id = $8
            WHERE s.status_name = $3
              AND ($8::INTEGER IS NULL OR assignee.id IS NOT NULL)
        )
        INSERT INTO tasks (
            title,
            description,
            status_id,
            priority_id,
            deadline,
            team_id,
            created_by,
            assigned_to
        )
        SELECT
            $1,
            $2,
            status_id,
            priority_id,
            $5,
            team_id,
            created_by,
            assigned_to
        FROM refs
        RETURNING id`,
		req.Title,
		req.Description,
		req.StatusName,
		req.PriorityID,
		req.Deadline,
		req.TeamID,
		req.CreatedBy,
		req.AssignedTo,
	).Scan(&taskID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Task{}, ErrRelatedEntityNotFound
		}
		return Task{}, fmt.Errorf("insert task: %w", err)
	}

	return DBTaskByID(ctx, taskID)
}

func DBUpdateTask(ctx context.Context, taskID int, req UpdateTaskRequest) (Task, error) {
	setParts := []string{}
	args := []any{}
	nextArg := 1

	addSet := func(column string, value any) {
		setParts = append(setParts, fmt.Sprintf("%s = $%d", column, nextArg))
		args = append(args, value)
		nextArg++
	}

	if req.Title != nil {
		addSet("title", strings.TrimSpace(*req.Title))
	}
	if req.Description != nil {
		addSet("description", *req.Description)
	}
	if req.StatusName != nil {
		statusID, err := DBStatusID(ctx, strings.TrimSpace(*req.StatusName))
		if err != nil {
			return Task{}, err
		}
		addSet("status_id", statusID)
	}
	if req.PriorityID != nil {
		ok, err := DBIDExists(ctx, "priorities", *req.PriorityID)
		if err != nil {
			return Task{}, err
		}
		if !ok {
			return Task{}, ErrRelatedEntityNotFound
		}
		addSet("priority_id", *req.PriorityID)
	}
	if req.Deadline != nil {
		addSet("deadline", req.Deadline)
	}
	if req.TeamID != nil {
		ok, err := DBIDExists(ctx, "teams", *req.TeamID)
		if err != nil {
			return Task{}, err
		}
		if !ok {
			return Task{}, ErrRelatedEntityNotFound
		}
		addSet("team_id", *req.TeamID)
	}
	if req.AssignedTo != nil {
		ok, err := DBIDExists(ctx, "users", *req.AssignedTo)
		if err != nil {
			return Task{}, err
		}
		if !ok {
			return Task{}, ErrRelatedEntityNotFound
		}
		addSet("assigned_to", *req.AssignedTo)
	}

	if len(setParts) == 0 {
		return Task{}, errors.New("empty update")
	}

	setParts = append(setParts, "updated_at = now()")
	args = append(args, taskID)
	query := fmt.Sprintf("UPDATE tasks SET %s WHERE id = $%d", strings.Join(setParts, ", "), nextArg)

	commandTag, err := Pool.Exec(ctx, query, args...)
	if err != nil {
		return Task{}, fmt.Errorf("update task: %w", err)
	}
	if commandTag.RowsAffected() == 0 {
		return Task{}, pgx.ErrNoRows
	}

	return DBTaskByID(ctx, taskID)
}

func DBUpdateTaskStatus(ctx context.Context, taskID int, statusName string) (Task, error) {
	statusID, err := DBStatusID(ctx, statusName)
	if err != nil {
		return Task{}, err
	}

	commandTag, err := Pool.Exec(ctx, `
        UPDATE tasks
        SET status_id = $1, updated_at = now()
        WHERE id = $2`, statusID, taskID)
	if err != nil {
		return Task{}, fmt.Errorf("update task status: %w", err)
	}
	if commandTag.RowsAffected() == 0 {
		return Task{}, pgx.ErrNoRows
	}

	return DBTaskByID(ctx, taskID)
}

func DBDeleteTask(ctx context.Context, taskID int) error {
	commandTag, err := Pool.Exec(ctx, "DELETE FROM tasks WHERE id = $1", taskID)
	if err != nil {
		return fmt.Errorf("delete task: %w", err)
	}
	if commandTag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func DBStatusID(ctx context.Context, statusName string) (int, error) {
	var statusID int
	err := Pool.QueryRow(ctx, "SELECT id FROM statuses WHERE status_name = $1", statusName).Scan(&statusID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, ErrRelatedEntityNotFound
		}
		return 0, fmt.Errorf("query status id: %w", err)
	}
	return statusID, nil
}

func DBIDExists(ctx context.Context, tableName string, id int) (bool, error) {
	allowedTables := map[string]bool{
		"priorities": true,
		"teams":      true,
		"users":      true,
	}
	if !allowedTables[tableName] {
		return false, fmt.Errorf("unsupported lookup table: %s", tableName)
	}

	var exists bool
	err := Pool.QueryRow(ctx, fmt.Sprintf("SELECT EXISTS (SELECT 1 FROM %s WHERE id = $1)", tableName), id).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("check %s id exists: %w", tableName, err)
	}
	return exists, nil
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

func DBCreateUser(ctx context.Context, req CreateUserRequest) (int, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}
	tx, err := Pool.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)

	var userID int
	err = tx.QueryRow(ctx,
		"INSERT INTO users (user_name, email) VALUES ($1, $2) RETURNING id",
		req.Username, req.Email).Scan(&userID)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return 0, fmt.Errorf("user already exists")
		}
		return 0, err
	}

	_, err = tx.Exec(ctx,
		"INSERT INTO login (user_name, email, password_hash, user_id) VALUES ($1, $2, $3, $4)",
		req.Username, req.Email, string(hashedPassword), userID)

	if err != nil {
		return 0, err
	}

	err = tx.Commit(ctx)
	return userID, err
}

func DBUpdateUser(ctx context.Context, userID int, req UpdateUserRequest) error {
	// 1. Защита от пустого запроса
	if req.Username == nil && req.Email == nil {
		return fmt.Errorf("payload must contain at least one field")
	}

	// 2. COALESCE оставит старое значение, если пришло NULL
	cmd, err := Pool.Exec(ctx, `
		UPDATE users SET
			user_name = COALESCE($2, user_name),
			email     = COALESCE($3, email),
			updated_at = NOW()
		WHERE id = $1
	`, userID, req.Username, req.Email)

	if err != nil {
		return err
	}

	// 3. Проверка, что пользователь действительно существует
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("user not found")
	}
	return nil
}

func DBDeleteUser(ctx context.Context, userID int) error {
	cmdTag, err := Pool.Exec(ctx,
		"DELETE FROM users WHERE id = $1",
		userID)

	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}
