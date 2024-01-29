package service

import (
	"cloud/external/ws"
	"cloud/grpc/users"
	"cloud/internal/explorer/http/request/share"
	"cloud/internal/explorer/model"
	"cloud/internal/explorer/repository"
	"context"
	"github.com/google/uuid"
	"github.com/patrickmn/go-cache"
	"log/slog"
	"time"
)

type ShareService struct {
	log   *slog.Logger
	cache *cache.Cache
	sr    repository.ShareRepositoryInterface
}

type ShareServiceInterface interface {
	Create(context.Context, *model.File, share.CreateRequest) (string, error)
	FindByFileID(context.Context, uint64) (*model.Share, error)
	FindByUserAndUUID(context.Context, string) (*model.Share, error)
	FindByUUID(string) (*model.Share, error)
	Update(context.Context, *model.Share, share.UpdateRequest) (*string, error)
	Delete(*model.Share) error
	UpdateDownloadCount(*model.Share) error
}

func NewShareService(
	log *slog.Logger,
	sr repository.ShareRepositoryInterface,
	cache *cache.Cache,
	ws *ws.WebSocketClient,
) *ShareService {
	return &ShareService{
		log:   log,
		sr:    sr,
		cache: cache,
		ws:    ws,
	}
}

func (s *ShareService) Create(ctx context.Context, file *model.File, req share.CreateRequest) (string, error) {
	const op = "ShareService.Create"

	log := s.log.With(
		slog.String("op", op),
	)

	user := ctx.Value("user").(*users.AuthorizeUserResponse)

	downloadLimit := int64(0)
	if req.DownloadLimit != nil {
		downloadLimit = *req.DownloadLimit
	}

	m := &model.Share{
		Uuid:          uuid.New().String(),
		FileID:        file.ID,
		UserID:        user.Id,
		Type:          req.Type,
		PinCode:       req.PinCode,
		DownloadLimit: downloadLimit,
		ExpiresAt:     req.ExpiresAt,
		CreatedAt:     time.Now(),
	}

	err := s.sr.Create(m)
	if err != nil {
		log.Error("failed to create share", slog.Any("err", err))

		return "", err
	}

	return m.Uuid, nil
}

func (s *ShareService) FindByFileID(ctx context.Context, fileID uint64) (*model.Share, error) {
	user := ctx.Value("user").(*users.AuthorizeUserResponse)

	return s.sr.FindByFileID(user.Id, fileID)
}

func (s *ShareService) FindByUserAndUUID(ctx context.Context, uuid string) (*model.Share, error) {
	user := ctx.Value("user").(*users.AuthorizeUserResponse)

	return s.sr.FindByUserIdAndUUID(user.Id, uuid)
}

func (s *ShareService) FindByUUID(uuid string) (*model.Share, error) {
	return s.sr.FindByUUID(uuid)
}

func (s *ShareService) Update(ctx context.Context, share *model.Share, req share.UpdateRequest) (*string, error) {
	const op = "ShareService.Update"

	log := s.log.With(
		slog.String("op", op),
	)

	now := time.Now()

	downloadLimit := int64(0)
	if req.DownloadLimit != nil {
		downloadLimit = *req.DownloadLimit
	}

	share.Uuid = uuid.New().String()
	share.Type = req.Type
	share.PinCode = req.PinCode
	share.DownloadLimit = downloadLimit
	share.ExpiresAt = req.ExpiresAt
	share.UpdatedAt = &now

	err := s.sr.Update(share)
	if err != nil {
		log.Error("failed to update share", slog.Any("err", err))

		return nil, err
	}

	return &share.Uuid, nil
}

func (s *ShareService) Delete(share *model.Share) error {
	const op = "ShareService.Delete"

	log := s.log.With(
		slog.String("op", op),
	)

	err := s.sr.Delete(share.Uuid)
	if err != nil {
		log.Error("failed to delete share", slog.Any("err", err))

		return err
	}

	return nil
}

func (s *ShareService) UpdateDownloadCount(share *model.Share) error {
	const op = "ShareService.UpdateDownloadCount"

	log := s.log.With(
		slog.String("op", op),
	)

	err := s.sr.UpdateDownloadCount(share.Uuid)
	if err != nil {
		log.Error("failed to update download count", slog.Any("err", err))

		return err
	}

	return nil
}
