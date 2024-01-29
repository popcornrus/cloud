package share

import (
	"cloud/internal/explorer/enum"
	"time"
)

type CreateRequest struct {
	FileUUID      string         `json:"file_uuid" validate:"required"`
	Type          enum.ShareType `json:"type" validate:"required"`
	DownloadLimit *int64         `json:"download_limit"`
	ExpiresAt     *time.Time     `json:"expires_at"`
	PinCode       *string        `json:"pin_code"`
}
