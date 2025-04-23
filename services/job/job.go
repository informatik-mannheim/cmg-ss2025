package models

import "time"

type JobStatus string

const (
	StatusQueued    JobStatus = "queued"
	StatusRunning   JobStatus = "running"
	StatusCompleted JobStatus = "completed"
	StatusFailed    JobStatus = "failed"
	StatusCancelled JobStatus = "cancelled"
)

type Job struct {
	ID             string         `json:"id"`                       // UUID
	UserID         string         `json:"userId"`                   // Referenz auf den Nutzer
	JobName        string         `json:"jobName"`                  // Freier Name für den Job
	Payload        string         `json:"payload"`                  // Eingabedaten oder Code
	Priority       string         `json:"priority"`                 // z.B. "low", "normal", "high"
	Status         JobStatus      `json:"status"`                   // Aktueller Status
	WorkerID       string         `json:"workerId,omitempty"`       // ID des ausführenden Workers
	Constraints    JobConstraints `json:"constraints"`              // Region, Laufzeit etc.
	Result         *string        `json:"result,omitempty"`         // Optionales Ergebnis (z. B. URL, Text)
	CreatedAt      time.Time      `json:"createdAt"`                // Zeitstempel der Erstellung
	StartedAt      *time.Time     `json:"startedAt,omitempty"`      // Startzeit
	CompletedAt    *time.Time     `json:"completedAt,omitempty"`    // Abschlusszeit
	EstimatedStart *time.Time     `json:"estimatedStart,omitempty"` // Optional geplant vom Scheduler
}

type JobConstraints struct {
	MaxRuntime       int      `json:"maxRuntime"`       // in Sekunden
	PreferredRegions []string `json:"preferredRegions"` // z. B. ["eu-central", "us-west"]
}
