package handler

import (
	"cloud/external/response"
	"cloud/external/ws"
	"cloud/grpc/users"
	"cloud/internal/explorer/enum"
	"cloud/internal/explorer/http/request/share"
	"cloud/internal/explorer/service"
	_struct "cloud/internal/explorer/struct"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
	"os"
)

type (
	ShareHandler struct {
		log       *slog.Logger
		validator *validator.Validate

		hs service.ShareServiceInterface
		fs service.FileServiceInterface
		ws *ws.WebSocketClient
	}

	ShareHandlerInterface interface {
		Create(w http.ResponseWriter, r *http.Request)
		Data(w http.ResponseWriter, r *http.Request)
		Show(w http.ResponseWriter, r *http.Request)
		Update(w http.ResponseWriter, r *http.Request)
		Delete(w http.ResponseWriter, r *http.Request)
		Download(w http.ResponseWriter, r *http.Request)
	}
)

func NewShareHandler(
	log *slog.Logger,
	hs service.ShareServiceInterface,
	fs service.FileServiceInterface,
	ws *ws.WebSocketClient,
) *ShareHandler {
	return &ShareHandler{
		log:       log,
		validator: validator.New(),
		hs:        hs,
		fs:        fs,
		ws:        ws,
	}
}

func (h *ShareHandler) Create(w http.ResponseWriter, r *http.Request) {
	const op = "ShareHandler.Create"

	log := h.log.With(
		slog.String("op", op),
	)

	var req share.CreateRequest

	if err := render.DecodeJSON(r.Body, &req); err != nil {
		log.Error("failed to decode request", err)

		response.Respond(w, response.Response{
			Status:  http.StatusBadRequest,
			Message: "bad request",
		})

		return
	}

	if err := h.validator.Struct(req); err != nil {
		log.Error("failed to validate request", err)

		response.Respond(w, response.Response{
			Status:  http.StatusUnprocessableEntity,
			Message: "unprocessable entity",
		})

		return
	}

	file, err := h.fs.FindByUUID(r.Context(), req.FileUUID)

	uuid, err := h.hs.Create(r.Context(), file, req)
	if err != nil {
		log.Error("failed to create share", err)

		response.Respond(w, response.Response{
			Status:  http.StatusInternalServerError,
			Message: "internal server error",
		})

		return
	}

	response.Respond(w, response.Response{
		Status: http.StatusOK,
		Data: map[string]string{
			"uuid": uuid,
		},
	})
}

func (h *ShareHandler) Update(w http.ResponseWriter, r *http.Request) {
	const op = "ShareHandler.Update"

	log := h.log.With(
		slog.String("op", op),
	)

	var req share.UpdateRequest

	if err := render.DecodeJSON(r.Body, &req); err != nil {
		log.Error("failed to decode request", err)

		response.Respond(w, response.Response{
			Status:  http.StatusBadRequest,
			Message: "bad request",
		})

		return
	}

	if err := h.validator.Struct(req); err != nil {
		log.Error("failed to validate request", err)

		response.Respond(w, response.Response{
			Status:  http.StatusUnprocessableEntity,
			Message: "unprocessable entity",
		})

		return
	}

	s, err := h.hs.FindByUserAndUUID(r.Context(), chi.URLParam(r, "uuid"))
	if err != nil {
		log.Error("failed to find share", err)

		response.Respond(w, response.Response{
			Status:  http.StatusNotFound,
			Message: "not found",
		})

		return
	}

	uuid, err := h.hs.Update(r.Context(), s, req)
	if err != nil {
		log.Error("failed to update share", err)

		response.Respond(w, response.Response{
			Status:  http.StatusInternalServerError,
			Message: "internal server error",
		})

		return
	}

	response.Respond(w, response.Response{
		Status: http.StatusOK,
		Data: map[string]string{
			"uuid": *uuid,
		},
	})
}

func (h *ShareHandler) Data(w http.ResponseWriter, r *http.Request) {
	const op = "ShareHandler.Data"

	log := h.log.With(
		slog.String("op", op),
	)

	file, err := h.fs.FindByUUID(r.Context(), chi.URLParam(r, "file"))
	if err != nil {
		log.Error("failed to find file", err)

		response.Respond(w, response.Response{
			Status:  http.StatusNotFound,
			Message: "not found",
		})

		return
	}

	s, err := h.hs.FindByFileID(r.Context(), file.ID)
	if err != nil {
		log.Error("failed to find share", err)

		response.Respond(w, response.Response{
			Status:  http.StatusNotFound,
			Message: "not found",
		})

		return
	}

	response.Respond(w, response.Response{
		Status: http.StatusOK,
		Data:   s,
	})
}

func (h *ShareHandler) Show(w http.ResponseWriter, r *http.Request) {
	const op = "ShareHandler.Show"

	log := h.log.With(
		slog.String("op", op),
	)

	s, err := h.hs.FindByUUID(chi.URLParam(r, "uuid"))
	if err != nil {
		log.Error("failed to find share", err)

		response.Respond(w, response.Response{
			Status:  http.StatusNotFound,
			Message: "not found",
		})

		return
	}

	if s.DownloadLimit != 0 && s.DownloadLimit == s.DownloadCount {
		response.Respond(w, response.Response{
			Status:  http.StatusGone,
			Message: "gone",
		})

		return
	}

	requiredPinCode := false
	if s.PinCode != nil {
		requiredPinCode = true
	}

	if chi.URLParam(r, "pin") != "" {
		hash := sha256.Sum256([]byte(*(s.PinCode)))

		if s.PinCode != nil && hex.EncodeToString(hash[:])[:32] != chi.URLParam(r, "pin") {
			response.Respond(w, response.Response{
				Status:  http.StatusForbidden,
				Message: "forbidden",
			})

			return
		}

		requiredPinCode = false
	}

	if requiredPinCode {
		result := _struct.PinResult{
			Uuid:      s.Uuid,
			PinCode:   requiredPinCode,
			CreatedAt: s.CreatedAt,
		}

		response.Respond(w, response.Response{
			Status: http.StatusOK,
			Data:   result,
		})
	} else {
		user, err := users.Get(log, s.UserID)
		if err != nil {
			log.Error("failed to get user", slog.Any("err", err))
		}

		file, err := h.fs.FindByID(r.Context(), s.FileID)
		if err != nil {
			log.Error("failed to find file", err)

			response.Respond(w, response.Response{
				Status:  http.StatusNotFound,
				Message: "not found",
			})

			return
		}

		if s.Type == enum.BurnType {
			ws.SendEvent(h.ws, ws.Socket{
				Channel: fmt.Sprintf("explorer.%s", user.Uuid),
				Event:   "send:share:burned",
				Data: map[string]string{
					"uuid": s.Uuid,
				},
			})
		}

		err = h.hs.UpdateDownloadCount(s)
		if err != nil {
			h.log.Error("failed to update download count", err)
		}

		if s.Type == enum.DownloadsLimitType {
			ws.SendEvent(h.ws, ws.Socket{
				Channel: fmt.Sprintf("explorer.%s", user.Uuid),
				Event:   "send:share:opened",
				Data: map[string]string{
					"uuid": s.Uuid,
				},
			})
		}

		if s.DownloadLimit != 0 && s.DownloadCount == s.DownloadLimit {
			ws.SendEvent(h.ws, ws.Socket{
				Channel: fmt.Sprintf("explorer.%s", user.Uuid),
				Event:   "send:share:reached_limit",
				Data: map[string]string{
					"uuid": s.Uuid,
				},
			})
		}

		response.Respond(w, response.Response{
			Status: http.StatusOK,
			Data: _struct.ShareResult{
				Uuid: s.Uuid,
				File: _struct.ShareFileResult{
					Uuid:    file.UUID,
					Name:    file.Name,
					Size:    file.Size,
					Type:    file.Type,
					Preview: file.Preview,
				},
			},
		})
	}

}

func (h *ShareHandler) Delete(w http.ResponseWriter, r *http.Request) {

}

func (h *ShareHandler) Download(w http.ResponseWriter, r *http.Request) {
	const op = "FileHandler.Download"

	log := h.log.With(
		slog.String("op", op),
	)

	file, err := h.fs.FindByUUID(r.Context(), chi.URLParam(r, "uuid"))
	if err != nil {
		log.Error("failed to find file by uuid", slog.Any("err", err))

		response.Respond(w, response.Response{
			Status: http.StatusNotFound,
		})

		return
	}

	path := fmt.Sprintf("%s/%s/%s", os.Getenv("SRV_PATH"), file.Path, file.Hash)

	f, err := os.Open(path)
	if err != nil {
		log.Error("failed to open video file", slog.Any("err", err))
		http.Error(w, fmt.Sprintf("Error opening video file: %s", err), http.StatusInternalServerError)
		return
	}

	defer f.Close()

	if file.IsVideo() {
		w.Header().Set("Content-Type", file.Type)
		w.Header().Set("Content-Length", fmt.Sprintf("%d", file.Size))
		w.Header().Set("Accept-Ranges", "bytes")

		_, err = f.Seek(0, 0)
		if err != nil {
			log.Error("failed to seek file", slog.Any("err", err))
			http.Error(w, fmt.Sprintf("Error seeking video file: %s", err), http.StatusInternalServerError)
			return
		}
	}

	http.ServeContent(w, r, file.Name, *file.UpdatedAt, f)
}
