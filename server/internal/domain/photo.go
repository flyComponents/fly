package domain

import "time"

type Photo struct {
	ID         string    `json:"id"`
	DeviceID   string    `json:"device_id"`
	UploadedAt time.Time `json:"uploaded_at"`
	Path       string    `json:"-"`
}
