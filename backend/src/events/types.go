package events

import (
	"time"

	"github.com/google/uuid"
)

type EventType string

const (
	EventTaskCreated         EventType = "task.created"
	EventTaskUpdated         EventType = "task.updated"
	EventTaskStatusChanged   EventType = "task.status_changed"
	EventTaskDeleted         EventType = "task.deleted"
	EventTeamCreated         EventType = "team.created"
	EventTeamSettingsUpdated EventType = "team.settings_updated"
)

type Event struct {
	EventID    string      `json:"event_id"`
	EventType  EventType   `json:"event_type"`
	OccurredAt time.Time   `json:"occurred_at"`
	Payload    interface{} `json:"payload"`
}

func NewEvent(eventType EventType, payload interface{}) Event {
	return Event{
		EventID:    uuid.New().String(),
		EventType:  eventType,
		OccurredAt: time.Now().UTC(),
		Payload:    payload,
	}
}

type TaskStatusPayload struct {
	TaskID    int    `json:"task_id"`
	TeamID    int    `json:"team_id"`
	OldStatus string `json:"old_status,omitempty"`
	NewStatus string `json:"new_status"`
}

type TaskPayload struct {
	TaskID  int         `json:"task_id"`
	TeamID  int         `json:"team_id"`
	Title   string      `json:"title,omitempty"`
	Status  string      `json:"status_name,omitempty"`
	Changes interface{} `json:"changes,omitempty"` // для updated
}

type TeamPayload struct {
	TeamID  int         `json:"team_id"`
	Name    string      `json:"name,omitempty"`
	Changes interface{} `json:"changes,omitempty"`
}
