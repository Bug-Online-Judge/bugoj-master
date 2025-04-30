package model

import "time"

type Problem struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Title       string    `gorm:"not null" json:"title"`
	Description string    `gorm:"type:text" json:"description"`
	TimeLimit   int       `json:"time_limit"`   // ms
	MemoryLimit int       `json:"memory_limit"` // MB
	DataVersion string    `json:"data_version"` // UUID
	CreatedAt   time.Time `json:"created_at"`
}
