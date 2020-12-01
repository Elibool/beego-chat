package models

import "time"

type Chats struct {
	Id          int64
	MemberA     int64
	MemberB     int64
	CreatedAt   time.Time `orm:"auto_now_add;type(datetime)"`
	UpdatedAt 	time.Time `orm:"auto_now;type(datetime)"`
}