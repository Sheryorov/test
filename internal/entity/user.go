package entity

import "time"

type User struct {
	ID        uint `gorm:"primaryKey"`
	FirstName string
	LastName  string
	Birthday  *time.Time
	Password  string
}
