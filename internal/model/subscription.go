package model

import (
    "time"
    "gorm.io/gorm"
)

type Subscription struct {
    ID         string         `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
    TargetURL  string         `gorm:"not null" json:"target_url"`
    Secret     string         `json:"secret,omitempty"`
    EventTypes []string       `gorm:"type:text[]" json:"event_types,omitempty"`
    CreatedAt  time.Time      `json:"created_at"`
    UpdatedAt  time.Time      `json:"updated_at"`
    DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}
