package UserLoginLogic

import "time"

type UserSession struct {
	ID            int       `json:"id"`
	Name          string    `json:"name"`
	PingCount     int       `json:"pingCount"`
	PingTime      time.Time `json:"pingTime"`
	ExpireSession time.Time `json:"expire"`
}
