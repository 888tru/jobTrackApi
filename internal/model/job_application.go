package model

import "time"

type ApplicationStatus string

const (
	StatusApplied     ApplicationStatus = "applied"
	StatusPhoneScreen ApplicationStatus = "phone_screen"
	StatusInterview   ApplicationStatus = "interview"
	StatusOffer       ApplicationStatus = "offer"
	StatusRejected    ApplicationStatus = "rejected"
	StatusAccepted    ApplicationStatus = "accepted"
	StatusWithdrawn   ApplicationStatus = "withdrawn"
)

type JobApplication struct {
	ID         uint              `gorm:"primaryKey" json:"id"`
	UserID     uint              `gorm:"not null;index" json:"-"`
	Company    string            `gorm:"not null" json:"company"`
	Position   string            `gorm:"not null" json:"position"`
	Status     ApplicationStatus `gorm:"not null;default:'applied'" json:"status"`
	AppliedAt  time.Time         `json:"applied_at"`
	Notes      string            `json:"notes"`
	URL        string            `json:"url"`
	CreatedAt  time.Time         `json:"created_at"`
	UpdatedAt  time.Time         `json:"updated_at"`
	Interviews []Interview       `gorm:"foreignKey:JobApplicationID" json:"interviews,omitempty"`
}
