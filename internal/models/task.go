package models

import (
	"time"
)

type Urgency int16

const (
	Low Urgency = iota
	Medium
	High
	Critical
)

type Task struct {
	Name        string
	Description string
	Urgency     Urgency
	DueDate     time.Time
}
