package API

import "time"

type Task struct {
	ID           int       `json:"id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	PriorityName string    `json:"priority_name"`
	StatusName   string    `json:"status_name"`
	Deadline     time.Time `json:"deadline"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	TeamName     string    `json:"team_name"`
	CreatedBy    string    `json:"created_by"`
	AssignedTo   string    `json:"assigned_to"`
}
