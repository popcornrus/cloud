package service

import (
	"cloud/external/ws"
	"cloud/internal/explorer/repository"
	"github.com/patrickmn/go-cache"
	"log/slog"
	"sync"
)

type ShareService struct {
	log   *slog.Logger
	cache *cache.Cache
	mu    sync.Mutex
	sr    repository.ShareRepositoryInterface

	ws *ws.WebSocketClient
}

type ShareServiceInterface interface {
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
