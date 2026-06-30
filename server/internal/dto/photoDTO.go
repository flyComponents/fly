package dto

type PhotoResponse struct {
	ID         string `json:"id"`
	DeviceID   string `json:"device_id"`
	UploadedAt string `json:"uploaded_at"`
	URL        string `json:"url"`
}
