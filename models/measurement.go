package models

import "time"

type Measurement struct {
	ID      int
	UserID  int
	Date    time.Time
	Weight  float64
	BodyFat float64
	Note    string
}
