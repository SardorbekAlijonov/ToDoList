package models

import (
	"time"
)

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:100;not null;" json:"name"`
	Email     string    `gorm:"size:100;not null;" json:"email"`
	Password  string    `gorm:"size:100;not null;" json:"password"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	Tasks     []Task    `gorm:"foreignkey:UserID" json:"tasks"`
}

type Task struct {
	ID          uint `gorm:"primaryKey"`
	Title       string
	Description string
	Status      string
	Priority    string
	Deadline    time.Time
	UserID      uint  `gorm:"not null"` // Foreign key for User
	User        User  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Tags        []Tag `gorm:"many2many:task_tags;"`
}
type Tag struct {
	ID    uint `gorm:"primaryKey"`
	Name  string
	Tasks []Task `gorm:"many2many:task_tags;"`
}
