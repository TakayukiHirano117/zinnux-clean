// Package model
package model

import "time"

type Project struct {
	ID          uint      `gorm:"primaryKey"`
	ProjectName string    `gorm:"not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
