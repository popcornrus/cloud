package service

import (
	"cloud/external/ws"
	"cloud/internal/explorer/repository"
	"github.com/patrickmn/go-cache"
	"log/slog"
	"sync"
)

type FolderService struct {
	log      *slog.Logger
	cache    *cache.Cache
	mu       sync.Mutex
	fileRepo repository.FolderRepositoryInterface

	ws *ws.WebSocketClient
}

type FolderServiceInterface interface {
}

func NewFolderService(
	log *slog.Logger,
	fr repository.FolderRepositoryInterface,
	cache *cache.Cache,
	ws *ws.WebSocketClient,
) *FolderService {
	return &FolderService{
		log:      log,
		fileRepo: fr,
		cache:    cache,
		ws:       ws,
	}
}
