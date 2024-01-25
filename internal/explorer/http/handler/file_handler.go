package handler

import (
	"bytes"
	"cloud/external/response"
	"cloud/grpc/users"
	"cloud/internal/explorer/enum"
	"cloud/internal/explorer/http/request/files"
	"cloud/internal/explorer/service"
	_struct "cloud/internal/explorer/struct"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"image/png"
	"log/slog"
	"net/http"
	"os"
	"time"
)

type (
	FileHandler struct {
		log       *slog.Logger
		validator *validator.Validate

		fs service.FileServiceInterface
	}

	FileHandlerInterface interface {
		List(http.ResponseWriter, *http.Request)
		Show(http.ResponseWriter, *http.Request)
		Data(http.ResponseWriter, *http.Request)
		Preview(http.ResponseWriter, *http.Request)
		Create(http.ResponseWriter, *http.Request)
		Update(http.ResponseWriter, *http.Request)
		Prepare(http.ResponseWriter, *http.Request)
		Download(http.ResponseWriter, *http.Request)
		Upload(http.ResponseWriter, *http.Request)
		Delete(http.ResponseWriter, *http.Request)
	}
)

func NewFileHandler(
	log *slog.Logger,
	fs service.FileServiceInterface,
) *FileHandler {
	return &FileHandler{
		log:       log,
		validator: validator.New(),
		fs:        fs,
	}
}

func (fh *FileHandler) List(w http.ResponseWriter, r *http.Request) {
	const op = "FileHandler.List"

	log := fh.log.With(
		slog.String("op", op),
	)

	user := r.Context().Value("user").(*users.AuthorizeUserResponse)
	if len(r.URL.Query().Get("search")) > 0 {
		f, err := fh.fs.Search(r.Context(), user, r.URL.Query().Get("search"))
		if err != nil {
			log.Error("failed to search files", slog.Any("err", err))
			return
		}

		response.Respond(w, response.Response{
			Status:  http.StatusOK,
			Message: "success",
			Data:    f,
		})

		return
	}

	f, err := fh.fs.List(r.Context(), user)
	if err != nil {
		log.Error("failed to list files", slog.Any("err", err))
		return
	}

	response.Respond(w, response.Response{
		Status:  http.StatusOK,
		Message: "success",
		Data:    f,
	})
}

func (fh *FileHandler) Show(w http.ResponseWriter, r *http.Request) {
	const op = "FileHandler.Show"

	log := fh.log.With(
		slog.String("op", op),
	)

	file, err := fh.fs.FindByUUID(r.Context(), chi.URLParam(r, "uuid"))
	if err != nil {
		log.Error("failed to find file by uuid", slog.Any("err", err))

		response.Respond(w, response.Response{
			Status: http.StatusNotFound,
		})

		return
	}

	path := fmt.Sprintf("%s/%s/%s", os.Getenv("SRV_PATH"), file.Path, file.Hash)
	if file.IsVideo() {
		path = fmt.Sprintf("%s/%s/%s.webm", os.Getenv("SRV_PATH"), file.Path, file.Hash)
	}

	f, err := os.Open(path)
	if err != nil {
		log.Error("failed to open video file", slog.Any("err", err))
		http.Error(w, fmt.Sprintf("Error opening video file: %s", err), http.StatusInternalServerError)
		return
	}

	defer f.Close()

	if file.IsVideo() {
		w.Header().Set("Content-Type", "video/webm")
		w.Header().Set("Content-Length", fmt.Sprintf("%d", file.Size))
		w.Header().Set("Accept-Ranges", "bytes")

		start, _, _, err := service.ParseRange(r.Header.Get("Range"), enum.WrappedContentRange)
		if err != nil {
			log.Error("failed to parse range", slog.Any("err", err))
			http.Error(w, fmt.Sprintf("Error parsing range: %s", err), http.StatusInternalServerError)
			return
		}

		_, err = f.Seek(int64(start), 0)
		if err != nil {
			log.Error("failed to seek file", slog.Any("err", err))
			http.Error(w, fmt.Sprintf("Error seeking video file: %s", err), http.StatusInternalServerError)
			return
		}
	}

	http.ServeContent(w, r, file.Name, *file.UpdatedAt, f)
}

func (fh *FileHandler) Data(w http.ResponseWriter, r *http.Request) {
	const op = "FileHandler.Data"

	log := fh.log.With(
		slog.String("op", op),
	)

	file, err := fh.fs.FindByUUID(r.Context(), chi.URLParam(r, "uuid"))
	if err != nil {
		log.Error("failed to find file by uuid", slog.Any("err", err))

		response.Respond(w, response.Response{
			Status: http.StatusNotFound,
		})

		return
	}

	response.Respond(w, response.Response{
		Status:  http.StatusOK,
		Message: "success",
		Data:    file,
	})
}

func (fh *FileHandler) Preview(w http.ResponseWriter, r *http.Request) {
	const op = "FileHandler.Preview"

	log := fh.log.With(
		slog.String("op", op),
	)

	file, err := fh.fs.FindByUUID(r.Context(), chi.URLParam(r, "uuid"))
	if err != nil {
		log.Error("failed to find file by uuid", slog.Any("err", err))

		response.Respond(w, response.Response{
			Status: http.StatusNotFound,
		})

		return
	}

	if file.Preview == nil {
		response.Respond(w, response.Response{
			Status: http.StatusAccepted,
		})

		return
	}

	image, err := fh.fs.Preview(_struct.PreviewProcessing{
		Width:  r.URL.Query().Get("w"),
		Height: r.URL.Query().Get("h"),
		Action: r.URL.Query().Get("a"),
	}, file)
	if err != nil {
		log.Error("failed to preview file", slog.Any("err", err))

		response.Respond(w, response.Response{
			Status: http.StatusInternalServerError,
		})

		return
	}

	buff := new(bytes.Buffer)
	err = png.Encode(buff, image)
	if err != nil {
		fmt.Println("failed to create buffer", err)
	}

	reader := bytes.NewReader(buff.Bytes())

	http.ServeContent(w, r, file.Name, time.Now(), reader)
	return
}

func (fh *FileHandler) Prepare(w http.ResponseWriter, r *http.Request) {
	const op = "FileHandler.Prepare"

	log := fh.log.With(
		slog.String("op", op),
	)

	var req files.PrepareRequest

	if err := render.DecodeJSON(r.Body, &req); err != nil {
		log.Error("failed to decode request", err)

		response.Respond(w, response.Response{
			Status:  http.StatusBadRequest,
			Message: "bad request",
		})

		return
	}

	if err := fh.validator.Struct(req); err != nil {
		log.Error("failed to validate request", slog.Any("err", err))
		return
	}

	result, err := fh.fs.Prepare(r.Context(), req)
	if err != nil {
		log.Error("failed to prepare file", slog.Any("err", err))
		return
	}

	response.Respond(w, response.Response{
		Status:  http.StatusOK,
		Message: "success",
		Data:    result,
	})
}

func (fh *FileHandler) Create(w http.ResponseWriter, r *http.Request) {

}

func (fh *FileHandler) Update(w http.ResponseWriter, r *http.Request) {
	const op = "FileHandler.Upload"

	log := fh.log.With(
		slog.String("op", op),
	)

	file, err := fh.fs.FindByUUID(r.Context(), chi.URLParam(r, "uuid"))
	if err != nil {
		log.Error("failed to find file by uuid", slog.Any("err", err))

		response.Respond(w, response.Response{
			Status: http.StatusNotFound,
		})

		return
	}

	var req files.UpdateRequest

	if err := render.DecodeJSON(r.Body, &req); err != nil {
		log.Error("failed to decode request", err)

		response.Respond(w, response.Response{
			Status:  http.StatusBadRequest,
			Message: "bad request",
		})

		return
	}

	if err := fh.validator.Struct(req); err != nil {
		log.Error("failed to validate request", slog.Any("err", err))

		response.Respond(w, response.Response{
			Status:  http.StatusUnprocessableEntity,
			Message: "bad request",
		})

		return
	}

	if err := fh.fs.Update(r.Context(), file, req); err != nil {
		log.Error("failed to update file", slog.Any("err", err))

		response.Respond(w, response.Response{
			Status:  http.StatusInternalServerError,
			Message: "internal server error",
		})

		return
	}

	response.Respond(w, response.Response{
		Status:  http.StatusOK,
		Message: "success",
	})
}

func (fh *FileHandler) Download(w http.ResponseWriter, r *http.Request) {
	const op = "FileHandler.Download"

	log := fh.log.With(
		slog.String("op", op),
	)

	file, err := fh.fs.FindByUUID(r.Context(), chi.URLParam(r, "uuid"))
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

func (fh *FileHandler) Upload(w http.ResponseWriter, r *http.Request) {
	const op = "FileHandler.Upload"

	log := fh.log.With(
		slog.String("op", op),
	)

	file, err := fh.fs.FindByUUID(r.Context(), chi.URLParam(r, "uuid"))
	if err != nil {
		log.Error("failed to find file by uuid", slog.Any("err", err))

		response.Respond(w, response.Response{
			Status: http.StatusNotFound,
		})

		return
	}

	if err := r.ParseMultipartForm(32 << 20); err != nil {
		log.Error("failed to parse multipart form", slog.Any("err", err))

		response.Respond(w, response.Response{
			Status: http.StatusBadRequest,
		})

		return
	}

	chunk, _, err := r.FormFile("chunk")
	if err != nil {
		log.Error("failed to get chunk from form", slog.Any("err", err))

		response.Respond(w, response.Response{
			Status: http.StatusBadRequest,
		})

		return
	}

	defer chunk.Close()

	contentRange := r.Header.Get("Content-Range")

	if err := fh.fs.Upload(
		r.Context(),
		contentRange,
		file,
		chunk,
	); err != nil {
		log.Error("failed to upload chunk", slog.Any("err", err))

		response.Respond(w, response.Response{
			Status: http.StatusInternalServerError,
		})

		return
	}

	response.Respond(w, response.Response{
		Status: http.StatusCreated,
	})
}

func (fh *FileHandler) Delete(w http.ResponseWriter, r *http.Request) {
	const op = "FileHandler.Delete"

	log := fh.log.With(
		slog.String("op", op),
	)

	file, err := fh.fs.FindByUUID(r.Context(), chi.URLParam(r, "uuid"))
	if err != nil {
		log.Error("failed to find file by uuid", slog.Any("err", err))

		response.Respond(w, response.Response{
			Status: http.StatusNotFound,
		})

		return
	}

	if err := fh.fs.Delete(r.Context(), file); err != nil {
		log.Error("failed to delete file", slog.Any("err", err))

		response.Respond(w, response.Response{
			Status: http.StatusInternalServerError,
		})

		return
	}

	response.Respond(w, response.Response{
		Status: http.StatusOK,
	})
}
