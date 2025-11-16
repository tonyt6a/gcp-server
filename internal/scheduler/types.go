package scheduler

import (
	"time"

	"github.com/google/uuid"
)

type JobStatus string

const (
	StatusPending   JobStatus = "PENDING"
	StatusRunning   JobStatus = "RUNNING"
	StatusCompleted JobStatus = "COMPLETED"
	StatusFailed    JobStatus = "FAILED"
)

type Job struct {
	ID        uuid.UUID `json:"id"`
	Status    JobStatus `json:"status"`
	Priority  int       `json:"priority"`
	Payload   string    `json:"payload"` // keep simple for v1
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}