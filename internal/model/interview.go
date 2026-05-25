package model

import "time"

type InterviewType string

const (
	InterviewPhone     InterviewType = "phone"
	InterviewTechnical InterviewType = "technical"
	InterviewHR        InterviewType = "hr"
	InterviewOnsite    InterviewType = "onsite"
	InterviewOther     InterviewType = "other"
)

type Interview struct {
	ID               uint          `gorm:"primaryKey" json:"id"`
	JobApplicationID uint          `gorm:"not null;index" json:"job_application_id"`
	ScheduledAt      time.Time     `json:"scheduled_at"`
	Type             InterviewType `json:"type"`
	Notes            string        `json:"notes"`
	Result           string        `json:"result"`
	CreatedAt        time.Time     `json:"created_at"`
	UpdatedAt        time.Time     `json:"updated_at"`
}
