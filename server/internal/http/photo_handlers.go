package http

import (
	"encoding/json"
	"fly_server/internal/dto"
	"fly_server/internal/service"
	"net/http"
	"path/filepath"
	"time"

	"github.com/go-chi/chi/v5"
)

type PhotoHandler struct {
	service *service.PhotoService
}

func NewPhotoHandler(s *service.PhotoService) *PhotoHandler {
	return &PhotoHandler{service: s}
}

// Upload godoc
// @Summary      Upload photo
// @Description  Upload photo from drone/device. Requires X-Device-ID header.
// @Tags         photos
// @Accept       multipart/form-data
// @Produce      json
// @Param        X-Device-ID header string true "Device ID"
// @Param        file formData file true "Photo file"
// @Success      200  {object}  dto.PhotoResponse
// @Failure      400  {string}  string "Bad request"
// @Failure      500  {string}  string "Internal error"
// @Router       /photos [post]
func (h *PhotoHandler) Upload(w http.ResponseWriter, r *http.Request) {
	deviceID := r.Header.Get("X-Device-ID")
	if deviceID == "" {
		http.Error(w, "missing device id", http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "file required", http.StatusBadRequest)
		return
	}
	defer file.Close()

	photo, err := h.service.Upload(r.Context(), deviceID, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := dto.PhotoResponse{
		ID:         photo.ID,
		DeviceID:   photo.DeviceID,
		UploadedAt: photo.UploadedAt.Format(time.RFC3339),
		URL:        "/files/" + photo.ID + ".jpg",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// GetPhoto godoc
// @Summary      Get photo metadata
// @Description  Returns metadata and URL for downloading the photo
// @Tags         photos
// @Produce      json
// @Param        id   path      string  true  "Photo ID"
// @Success      200  {object}  dto.PhotoResponse
// @Failure      404  {string}  string "Not found"
// @Failure      500  {string}  string "Internal error"
// @Router       /photos/{id} [get]
func (h *PhotoHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	photo, err := h.service.Get(r.Context(), id)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	resp := dto.PhotoResponse{
		ID:         photo.ID,
		DeviceID:   photo.DeviceID,
		UploadedAt: photo.UploadedAt.Format(time.RFC3339),
		URL:        "/files/" + photo.ID + ".jpg",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// ServeFile godoc
// @Summary      Download file
// @Description  Serves uploaded file by filename
// @Tags         files
// @Produce      application/octet-stream
// @Param        filename   path      string  true  "File name (e.g. 123.jpg)"
// @Success      200  {string}  binary  "File content"
// @Failure      404  {string}  string  "Not found"
// @Router       /files/{filename} [get]
func (h *PhotoHandler) ServeFile(w http.ResponseWriter, r *http.Request) {
	filename := filepath.Base(chi.URLParam(r, "filename"))

	http.ServeFile(w, r, "./uploads/"+filename)
}

// GetByOwnerId godoc
// @Summary      List photos by owner/device
// @Description  Returns a list of photos uploaded by a specific device (DeviceID)
// @Tags         photos
// @Produce      json
// @Param        id   path      string  true  "Device ID"
// @Success      200  {array}   dto.PhotoResponse
// @Failure      404  {string}  string "Not found"
// @Failure      500  {string}  string "Internal error"
// @Router       /photos/owner/{id} [get]
func (h *PhotoHandler) GetByOwnerId(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	photos, err := h.service.GetByOwnerId(r.Context(), id)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	resp := []dto.PhotoResponse{}
	for _, photo := range photos {
		resp = append(resp, dto.PhotoResponse{
			ID:         photo.ID,
			DeviceID:   photo.DeviceID,
			UploadedAt: photo.UploadedAt.Format(time.RFC3339),
			URL:        "/files/" + photo.ID + ".jpg",
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
