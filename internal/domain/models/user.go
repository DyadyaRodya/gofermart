package models

import "time"

type UserInfo struct {
	UUID         string
	Login        string
	CreatedAt    time.Time
	PasswordHash string
	PasswordSalt string
}
