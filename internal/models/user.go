package models

import (
	"strings"
	"time"
)

type User struct {
	ID        int       `gorm:"serial;primaryKey"`
	Name      string    `gorm:"size:100;"`
	DeletedAt time.Time `gorm:"default:null;"`
}

func NewUser(name string) *User {
	return &User{
		Name: name,
	}
}

func (u *User) Validate(create bool) bool {
	u.Name = strings.TrimSpace(u.Name)
	return len([]rune(u.Name)) > 0 && (create || u.ID > 0)
}
