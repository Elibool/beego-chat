package models

import "time"

type Messages struct {
	Id        int64     `json:"id"`
	ChatId    int64     `json:"chat_id"`
	MemberId  int64     `json:"member_id"`
	Msg       string    `json:"msg"`
	IfRead    int64     `json:"if_read"`
	CreatedAt time.Time `json:"created_at" orm:"auto_now_add;type(datetime)"`
	ReadAt    time.Time `json:"read_at" orm:"auto_now;type(datetime)"`
}
