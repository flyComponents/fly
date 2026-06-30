package domain

import "time"

type Agent struct {
	ID       string
	LastSeen time.Time
}
