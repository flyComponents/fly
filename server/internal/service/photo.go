package service

import (
	"context"
	"fly_server/internal/domain"
	"fly_server/internal/storage"
	"io"
	"time"
)

type PhotoService struct {
	repo storage.PhotoRepository
	fs   storage.FileStorage
}

func NewPhotoService(r storage.PhotoRepository, fs storage.FileStorage) *PhotoService {
	return &PhotoService{repo: r, fs: fs}
}

func (s *PhotoService) Upload(ctx context.Context, deviceID string, r io.Reader) (domain.Photo, error) {
	id := deviceID + "__" + time.Now().Format("2006-01-02_15-04-05")

	path, err := s.fs.Save(ctx, id, r)
	if err != nil {
		return domain.Photo{}, err
	}

	photo := domain.Photo{
		ID:         id,
		DeviceID:   deviceID,
		UploadedAt: time.Now().UTC(),
		Path:       path,
	}

	if err := s.repo.Save(ctx, photo); err != nil {
		return photo, err
	}

	return photo, nil
}

func (s *PhotoService) Get(ctx context.Context, ID string) (*domain.Photo, error) {
	photo, err := s.repo.Get(ctx, ID)
	if err != nil {
		return &domain.Photo{}, err
	}

	return photo, nil

}

func (s *PhotoService) GetByOwnerId(ctx context.Context, id string) ([]*domain.Photo, error) {
	photos, err := s.repo.GetByOwnerId(ctx, id)
	if err != nil {
		return []*domain.Photo{}, err
	}

	return photos, nil

}
