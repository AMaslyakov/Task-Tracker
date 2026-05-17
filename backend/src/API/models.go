package API

import "time"

type Task struct {
	ID           int        `json:"id"`
	Title        string     `json:"title"`
	Description  string     `json:"description"`
	PriorityName string     `json:"priority_name"`
	StatusName   string     `json:"status_name"`
	Deadline     *time.Time `json:"deadline"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	TeamID       int        `json:"team_id"`
	TeamName     string     `json:"team_name"`
	CreatedBy    string     `json:"created_by"`
	AssignedTo   string     `json:"assigned_to"`
}

type CreateTaskRequest struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	StatusName  string     `json:"status_name"`
	PriorityID  int        `json:"priority_id"`
	Deadline    *time.Time `json:"deadline"`
	TeamID      int        `json:"team_id"`
	CreatedBy   int        `json:"created_by"`
	AssignedTo  *int       `json:"assigned_to"`
}

type UpdateTaskRequest struct {
	Title       *string    `json:"title"`
	Description *string    `json:"description"`
	StatusName  *string    `json:"status_name"`
	PriorityID  *int       `json:"priority_id"`
	Deadline    *time.Time `json:"deadline"`
	TeamID      *int       `json:"team_id"`
	AssignedTo  *int       `json:"assigned_to"`
}

type UpdateTaskStatusRequest struct {
	StatusName string `json:"status_name"`
}

type Team struct {
	ID              int      `json:"id"`
	Name            string   `json:"name"`
	Description     string   `json:"description"`
	ConfigDashboard []string `json:"config_dashboard"`
	CreatedAt       string   `json:"created_at"`
	UpdatedAt       string   `json:"updated_at"`
	Tasks           []string `json:"tasks"`
	Members         []string `json:"members"`
}
