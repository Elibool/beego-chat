package models

import (
	"time"
)

type Members struct {
	Id          int64     `json:"id"`
	Name        string    `json:"name"`
	LoginIp     string    `json:"login_ip"`
	IsOnline    int       `json:"is_online"`
	CreatedAt   time.Time `json:"created_at" orm:"auto_now_add;type(datetime)"`
	LastLoginAt time.Time `json:"last_login_at" orm:"auto_now;type(datetime)"`
}

