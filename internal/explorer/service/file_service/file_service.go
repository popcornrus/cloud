package file_service

import (
	"cloud/external/ws"
	"cloud/grpc/users"
	"cloud/internal/explorer/enum"
	"cloud/internal/explorer/http/request/files"
	"cloud/internal/explorer/model"
	"cloud/internal/explorer/repository"
	_struct "cloud/internal/explorer/struct"
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/anthonynsimon/bild/imgio"
	"github.com/anthonynsimon/bild/transform"
	"github.com/google/uuid"
	"github.com/patrickmn/go-cache"
	"image"
	"io"
	"log/slog"
	"math/big"
	"mime/multipart"
	"os"
	"regexp"
	"strconv"
	"sync"
	"time"
)

type FileService struct {
	log      *slog.Logger
	cache    *cache.Cache
	mu       sync.Mutex
	fileRepo repository.FileRepositoryInterface

	ws *ws.WebSocketClient
}

type FileServiceInterface interface {
	List(context.Context, *users.AuthorizeUserResponse) ([]*model.File, error)
	FindByUUID(context.Context, string) (*model.File, error)
	Preview(_struct.PreviewProcessing, *model.File) (image.Image, error)
	Update(context.Context, *model.File, files.UpdateRequest) error
	Prepare(context.Context, files.PrepareRequest) (*_struct.PrepareResult, error)
	Upload(context.Context, string, *model.File, multipart.File) error
	Delete(context.Context, *model.File) error

	CollectFile(context.Context, *model.File) error
}

func NewFileService(
	log *slog.Logger,
	fr repository.FileRepositoryInterface,
	cache *cache.Cache,
	ws *ws.WebSocketClient,
) *FileService {
	return &FileService{
		log:      log,
		fileRepo: fr,
		cache:    cache,
		ws:       ws,
	}
}

const ChunkSize = 1024 * 50

func (fs *FileService) List(ctx context.Context, user *users.AuthorizeUserResponse) ([]*model.File, error) {
	const op = "FileService.List"

	log := fs.log.With(
		slog.String("op", op),
	)

	f, err := fs.fileRepo.GetByUserID(ctx, uint64(user.Id))

	if err != nil {
		log.Error("failed to get files", slog.Any("err", err))
		return nil, err
	}

	return f, nil
}

func (fs *FileService) Preview(proc _struct.PreviewProcessing, file *model.File) (image.Image, error) {
	const op = "FileService.Preview"

	log := fs.log.With(
		slog.String("op", op),
	)

	var width int
	var height int
	var err error

	if proc.Width != "" {
		width, err = strconv.Atoi(proc.Width)
		if err != nil {
			log.Error("failed to parse width", slog.Any("err", err))
			return nil, err
		}
	} else {
		return nil, errors.New("width is required")
	}

	if proc.Height != "" {
		height, err = strconv.Atoi(proc.Height)
		if err != nil {
			log.Error("failed to parse height", slog.Any("err", err))
			return nil, err
		}
	} else {
		return nil, errors.New("height is required")
	}

	cacheHash := _hash(fmt.Sprintf("%s-%s-%s", file.Hash, proc.Width, proc.Height))
	cachePath := fmt.Sprintf("%s/%s/.cache/%s/%s", os.Getenv("SRV_PATH"), file.Path, file.Hash, cacheHash)
	cacheFile := fmt.Sprintf("%s/%s", cachePath, cacheHash)

	if _, err := os.Stat(cacheFile); err == nil {
		img, err := imgio.Open(cacheFile)
		if err != nil {
			log.Error("failed to open image", slog.Any("err", err), slog.Any("path", cachePath))
			return nil, err
		}

		return img, nil
	}

	path := fmt.Sprintf("%s/%s/.preview/%s.jpg", os.Getenv("SRV_PATH"), file.Path, *file.Preview)

	img, err := imgio.Open(path)
	if err != nil {
		log.Error("failed to open image", slog.Any("err", err))
		return nil, err
	}

	var result *image.RGBA

	if proc.Action != "" {
		switch proc.Action {
		case "crop":
			result = transform.Crop(img, image.Rect(0, 0, width, height))
			break
		case "resize":
			result = transform.Resize(img, width, height, transform.Lanczos)
			break
		default:
			result = transform.Resize(img, width, height, transform.Lanczos)
			break
		}
	}

	if err := os.MkdirAll(cachePath, 0755); err != nil {
		log.Error("failed to create cache directory", slog.Any("err", err))
		return nil, err
	}

	if err := imgio.Save(cacheFile, result, imgio.JPEGEncoder(100)); err != nil {
		log.Error("failed to save cache file", slog.Any("err", err))
		return nil, err
	}

	return imgio.Open(cacheFile)
}

func (fs *FileService) Update(ctx context.Context, file *model.File, req files.UpdateRequest) error {
	if req.Name != "" {
		file.Name = req.Name
	}

	if err := fs.fileRepo.Update(ctx, file); err != nil {
		return err
	}

	return nil
}

func (fs *FileService) Prepare(ctx context.Context, req files.PrepareRequest) (*_struct.PrepareResult, error) {
	user := ctx.Value("user").(*users.AuthorizeUserResponse)

	fileHash := _hash(fmt.Sprintf("%s-%s-%d-%s", user.Uuid, req.Name, req.Size, time.Now()))

	file := &model.File{
		UserID: uint64(user.Id),
		UUID:   uuid.New().String(),
		Name:   req.Name,
		Path:   _createPathFromHash(new(big.Int).SetBytes([]byte(fileHash)).String()),
		Hash:   fileHash,
		State:  enum.FileStatePending,
		Size:   req.Size,
		Type:   req.Type,
	}

	var err error

	file.ID, err = fs.fileRepo.Create(ctx, file)
	if err != nil {
		return nil, err
	}

	fs.cache.Set(fmt.Sprintf("file-%s", file.UUID), file, 5*time.Minute)

	ws.SendEvent(fs.ws, ws.Socket{
		Channel: fmt.Sprintf("files.%s", user.Uuid),
		Event:   "send:file:created",
		Data: map[string]any{
			"uuid":  file.UUID,
			"state": file.State,
		},
	})

	return &_struct.PrepareResult{
		Url:       fmt.Sprintf("/%s/upload", file.UUID),
		ChunkSize: ChunkSize,
	}, nil
}

func (fs *FileService) Upload(
	ctx context.Context,
	contentRange string,
	file *model.File,
	f multipart.File,
) error {
	start, _, _, err := ParseRange(contentRange, enum.FullContentRange)
	if err != nil {
		return err
	}

	if start == 0 {
		file.State = enum.FileStateUploading

		if err := fs.fileRepo.Update(context.Background(), file); err != nil {
			return err
		}
	}

	chunkHash := _hash(fmt.Sprintf("%s-%s", file.Hash, contentRange))

	if err := os.MkdirAll(fmt.Sprintf("%s/%s", os.Getenv("SRV_PATH"), file.Path), 0755); err != nil {
		return err
	}

	chunkPath := fmt.Sprintf("%s/%s/%s", os.Getenv("SRV_PATH"), file.Path, chunkHash)
	out, err := os.Create(chunkPath)
	if err != nil {
		return err
	}

	defer out.Close()

	if _, err := io.Copy(out, f); err != nil {
		return err
	}

	if err := fs.fileRepo.StoreFileChunk(context.Background(), &model.FileChunk{
		FileID: file.ID,
		Chunk:  start,
		Hash:   chunkHash,
	}); err != nil {
		return err
	}

	if fs.FileIsUploaded(file) {
		file.State = enum.FileStateCollecting

		if err := fs.fileRepo.Update(context.Background(), file); err != nil {
			return err
		}

		user := ctx.Value("user").(*users.AuthorizeUserResponse)
		ws.SendEvent(fs.ws, ws.Socket{
			Channel: fmt.Sprintf("files.%s", user.Uuid),
			Event:   "send:file:update",
			Data: map[string]any{
				"uuid":  file.UUID,
				"state": file.State,
			},
		})

		go func(ctx context.Context, file *model.File) {
			err := fs.CollectFile(ctx, file)
			if err != nil {
				fmt.Println(err)
				return
			}
		}(ctx, file)
	}

	return nil
}

func (fs *FileService) Delete(ctx context.Context, file *model.File) error {
	const op = "FileService.Delete"

	log := fs.log.With(
		slog.String("op", op),
	)

	if err := fs.fileRepo.Delete(ctx, file); err != nil {
		log.Error("failed to delete file", slog.Any("err", err))
		return err
	}

	go func() {
		filePath := fmt.Sprintf("%s/%s/%s", os.Getenv("SRV_PATH"), file.Path, file.Hash)
		if _, err := os.Stat(filePath); err == nil {
			if err := os.Remove(filePath); err != nil {
				log.Error("failed to remove file", slog.Any("err", err))
				return
			}
		}

		previewPath := fmt.Sprintf("%s/%s/.preview", os.Getenv("SRV_PATH"), file.Path)
		previewFile := fmt.Sprintf("%s/%s.jpg", previewPath, file.Hash)
		if _, err := os.Stat(previewFile); err == nil {
			if err := os.Remove(previewFile); err != nil {
				log.Error("failed to remove preview file", slog.Any("err", err))
				return
			}
		}

		cachePath := fmt.Sprintf("%s/%s/.cache", os.Getenv("SRV_PATH"), file.Path)
		cacheFile := fmt.Sprintf("%s/%s", cachePath, file.Hash)
		if _, err := os.Stat(cacheFile); err == nil {
			if err := os.RemoveAll(cacheFile); err != nil {
				log.Error("failed to remove cache file", slog.Any("err", err))
				return
			}
		}

		if file.IsVideo() {
			webmPath := fmt.Sprintf("%s/%s/%s.webm", os.Getenv("SRV_PATH"), file.Path, file.Hash)
			if _, err := os.Stat(webmPath); err == nil {
				if err := os.Remove(webmPath); err != nil {
					log.Error("failed to remove webm file", slog.Any("err", err))
					return
				}
			}
		}
	}()

	return nil
}

func (fs *FileService) FindByUUID(ctx context.Context, uuid string) (*model.File, error) {
	cachedFile, ok := fs.cache.Get(fmt.Sprintf("file-%s", uuid))
	if ok {
		return cachedFile.(*model.File), nil
	}

	file, err := fs.fileRepo.FindByUUID(ctx, uuid)

	if err != nil {
		return nil, err
	}

	fs.cache.Set(fmt.Sprintf("file-%s", uuid), file, 5*time.Minute)

	return file, nil
}

func (fs *FileService) CollectFile(ctx context.Context, file *model.File) error {
	const op = "FileService.CollectFile"

	log := fs.log.With(
		slog.String("op", op),
	)

	chunks, err := fs.fileRepo.GetFileChunks(context.Background(), file.ID)
	if err != nil {
		log.Error("failed to get file chunks", slog.Any("err", err))
		return err
	}

	if len(chunks) == 0 {
		log.Error("no chunks found")
		return errors.New("no chunks found")
	}

	if err := os.MkdirAll(fmt.Sprintf("%s/%s", os.Getenv("SRV_PATH"), file.Path), 0755); err != nil {
		log.Error("failed to create directory", slog.Any("err", err))
		return err
	}

	filePath := fmt.Sprintf("%s/%s/%s", os.Getenv("SRV_PATH"), file.Path, file.Hash)
	out, err := os.Create(filePath)
	if err != nil {
		log.Error("failed to create file", slog.Any("err", err))
		return err
	}

	defer out.Close()

	for _, chunk := range chunks {
		chunkPath := fmt.Sprintf("%s/%s/%s", os.Getenv("SRV_PATH"), file.Path, chunk.Hash)

		f, err := os.Open(chunkPath)
		if err != nil {
			log.Error("failed to open chunk", slog.Any("err", err))
			return err
		}

		defer f.Close()

		if _, err := io.Copy(out, f); err != nil {
			log.Error("failed to copy chunk", slog.Any("err", err))
			return err
		}

		if err := os.Remove(chunkPath); err != nil {
			log.Error("failed to remove chunk", slog.Any("err", err))
			return err
		}
	}

	err = fs.fileRepo.RemoveFileChunks(context.Background(), file.ID)
	if err != nil {
		log.Error("failed to remove file chunks", slog.Any("err", err))
		return err
	}

	file.Preview, err = fs.CreatePreview(ctx, file)
	if err != nil {
		log.Error("failed to create preview", slog.Any("err", err))
		return err
	}

	user := ctx.Value("user").(*users.AuthorizeUserResponse)

	if file.IsVideo() {
		go func() {
			file.State = enum.FileStateConverting

			if err := fs.fileRepo.Update(context.Background(), file); err != nil {
				log.Error("failed to update file", slog.Any("err", err))
				return
			}

			ws.SendEvent(fs.ws, ws.Socket{
				Channel: fmt.Sprintf("files.%s", user.Uuid),
				Event:   "send:file:update",
				Data: map[string]any{
					"uuid":  file.UUID,
					"state": file.State,
				},
			})

			err := file.ConvertToWebM()
			if err != nil {
				log.Error("failed to convert to webm", slog.Any("err", err))
				return
			}

			file.State = enum.FileStateDone

			if err := fs.fileRepo.Update(context.Background(), file); err != nil {
				log.Error("failed to update file", slog.Any("err", err))
				return
			}

			ws.SendEvent(fs.ws, ws.Socket{
				Channel: fmt.Sprintf("files.%s", user.Uuid),
				Event:   "send:file:update",
				Data: map[string]any{
					"uuid":  file.UUID,
					"state": file.State,
				},
			})
		}()
	} else {
		file.State = enum.FileStateDone

		if err := fs.fileRepo.Update(context.Background(), file); err != nil {
			log.Error("failed to update file", slog.Any("err", err))
			return err
		}

		ws.SendEvent(fs.ws, ws.Socket{
			Channel: fmt.Sprintf("files.%s", user.Uuid),
			Event:   "send:file:update",
			Data: map[string]any{
				"uuid":  file.UUID,
				"state": file.State,
			},
		})
	}

	return nil
}

func (fs *FileService) CreatePreview(ctx context.Context, file *model.File) (*string, error) {
	const op = "FileService.CreatePreview"

	log := fs.log.With(
		slog.String("op", op),
	)

	var previewImage *string
	var err error

	if file.IsImage() {
		if previewImage, err = file.CreateImagePreview(); err != nil {
			log.Error("failed to create image preview", slog.Any("err", err))
			return nil, err
		}
	}

	if file.IsVideo() {
		if previewImage, err = file.CreateVideoPreview(); err != nil {
			log.Error("failed to create video preview", slog.Any("err", err))
			return nil, err
		}
	}

	user := ctx.Value("user").(*users.AuthorizeUserResponse)
	ws.SendEvent(fs.ws, ws.Socket{
		Channel: fmt.Sprintf("files.%s", user.Uuid),
		Event:   "send:file:preview",
		Data: map[string]any{
			"uuid": file.UUID,
		},
	})

	return previewImage, nil
}

func _hash(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func _createPathFromHash(hash string) string {
	return fmt.Sprintf("%s/%s/%s", hash[0:2], hash[2:4], hash[4:6])
}

func ParseRange(input string, rangeType enum.ContentRange) (int, int, int, error) {
	var re *regexp.Regexp

	re = regexp.MustCompile(string(rangeType))
	matches := re.FindStringSubmatch(input)

	if len(matches) < 2 {
		return 0, 0, 0, fmt.Errorf("Invalid range format")
	}

	start, end, total := 0, 0, 0

	start, err := strconv.Atoi(matches[1])
	if err != nil {
		return 0, 0, 0, err
	}

	if len(matches) > 3 {
		end, err = strconv.Atoi(matches[2])
		if err != nil {
			return 0, 0, 0, err
		}
	}

	if len(matches) > 4 {
		total, err = strconv.Atoi(matches[3])

		if err != nil {
			return 0, 0, 0, err
		}
	}

	return start, end, total, nil
}

var uploadedChunks = make(map[string]int)

func (fs *FileService) FileIsUploaded(file *model.File) bool {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	chunks, ok := uploadedChunks[file.UUID]

	if !ok {
		countOfChunks := int(file.Size / ChunkSize)
		if countOfChunks == 0 {
			return true
		}

		uploadedChunks[file.UUID] = countOfChunks
		return false
	}

	uploadedChunks[file.UUID] = chunks - 1

	if uploadedChunks[file.UUID] == 0 {
		delete(uploadedChunks, file.UUID)
		return true
	}

	return false
}
