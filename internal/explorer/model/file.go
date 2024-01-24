package model

import (
	"cloud/internal/explorer/enum"
	"fmt"
	"github.com/anthonynsimon/bild/imgio"
	"github.com/anthonynsimon/bild/transform"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"os"
	"os/exec"
	"time"
)

type (
	File struct {
		ID        uint64         `json:"-"`
		UserID    uint64         `json:"-"`
		UUID      string         `json:"uuid"`
		Name      string         `json:"name"`
		Path      string         `json:"path"`
		Hash      string         `json:"hash"`
		State     enum.FileState `json:"state"`
		Size      int64          `json:"size"`
		Type      string         `json:"type"`
		Preview   *string        `json:"preview"`
		CreatedAt time.Time      `json:"created_at"`
		UpdatedAt *time.Time     `json:"updated_at"`
	}

	FileChunk struct {
		ID        primitive.ObjectID `json:"-" bson:"_id,omitempty"`
		FileID    uint64             `json:"-" bson:"file_id"`
		Chunk     int                `json:"chunk" bson:"chunk"`
		Hash      string             `json:"hash" bson:"hash"`
		CreatedAt time.Time          `json:"created_at" bson:"created_at"`
		UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
	}
)

func (f *File) IsImage() bool {
	types := []string{
		"image/png",
		"image/jpg",
		"image/webp",
		"image/jpeg",
	}

	for _, t := range types {
		if t == f.Type {
			return true
		}
	}

	return false
}

func (f *File) IsVideo() bool {
	types := []string{
		"video/mp4",
		"video/ogg",
		"video/webm",
	}

	for _, t := range types {
		if t == f.Type {
			return true
		}
	}

	return false
}

func (f *File) CreateImagePreview() (*string, error) {
	previewPath := fmt.Sprintf("%s/%s/.preview", os.Getenv("SRV_PATH"), f.Path)
	previewFile := fmt.Sprintf("%s/%s.jpg", previewPath, f.Hash)

	file := fmt.Sprintf("%s/%s/%s", os.Getenv("SRV_PATH"), f.Path, f.Hash)

	if _, err := os.Stat(file); err != nil {
		return nil, err
	}

	if err := os.MkdirAll(previewPath, 0755); err != nil {
		return nil, err
	}

	img, err := imgio.Open(file)
	if err != nil {
		return nil, err
	}

	result := transform.Resize(img, 1280, 720, transform.Lanczos)

	if err := imgio.Save(previewFile, result, imgio.JPEGEncoder(100)); err != nil {
		return nil, err
	}

	return &f.Hash, nil
}

func (f *File) CreateVideoPreview() (*string, error) {
	previewPath := fmt.Sprintf("%s/%s/.preview", os.Getenv("SRV_PATH"), f.Path)
	previewFile := fmt.Sprintf("%s/%s.jpg", previewPath, f.Hash)

	file := fmt.Sprintf("%s/%s/%s", os.Getenv("SRV_PATH"), f.Path, f.Hash)

	if _, err := os.Stat(file); err != nil {
		return nil, err
	}

	if err := os.MkdirAll(previewPath, 0755); err != nil {
		return nil, err
	}

	cmd := exec.Command("ffmpeg", "-i", file, "-ss", "00:00:05", "-vframes:v", "1", previewFile)
	//cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("error creating thumbnail: %v", err)
	}

	return &f.Hash, nil
}

type VideoProcessing interface {
	ConvertToWebM() error
	ConvertToMP4() error
}

func (f *File) ConvertToWebM() error {
	file := fmt.Sprintf("%s/%s/%s", os.Getenv("SRV_PATH"), f.Path, f.Hash)
	webm := fmt.Sprintf("%s/%s/%s.webm", os.Getenv("SRV_PATH"), f.Path, f.Hash)

	cmd := exec.Command("ffmpeg", "-i", file, "-c:v", "libvpx-vp9", "-c:a", "libopus", webm)
	//cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("error converting to webm: %v", err)
	}

	return nil
}
