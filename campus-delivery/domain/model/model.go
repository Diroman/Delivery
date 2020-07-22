package model

import "time"

type User struct {
	Id           int64
	NickName     string
	FirstName    string
	SecondName   string
	Latitude     float32
	Longitude    float32
	Rating       float32
	Notification bool
}

type Courier struct {
	User        User
	Shop        string
	Description string
	TimeFrom    time.Time
	TimeTo      time.Time
	Link        string
	ChatId      int64
}

type Code struct {
	Code  string
	Error bool
}

type UserRequest struct {
	Id       int64
	TimeFrom time.Time
	TimeTo   time.Time
}

type Rating struct {
	Id     int64
	Rating float32
}

type ChanMessage struct {
	Ids     []int64
	Courier Courier
}
