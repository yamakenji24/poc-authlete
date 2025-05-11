package entity

import "time"

// User はユーザー情報を保持する構造体です
type User struct {
	ID        string
	Username  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// UserRepository はユーザー情報を管理するリポジトリのインターフェースです
type UserRepository interface {
	FindByID(id string) (*User, error)
	FindByUsername(username string) (*User, error)
	Save(user *User) error
	Delete(id string) error
}
