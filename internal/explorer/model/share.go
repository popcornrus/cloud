package model

import (
	"cloud/internal/explorer/enum"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Share struct {
	ID            primitive.ObjectID `json:"-" bson:"_id,omitempty"`
	Uuid          string             `json:"uuid" bson:"uuid"`
	FileID        uint64             `json:"-" bson:"file_id"`
	UserID        int64              `json:"-" bson:"user_id"`
	Type          enum.ShareType     `json:"type" bson:"type"`
	PinCode       *string            `json:"pin_code" bson:"pin_code"`
	DownloadLimit int64              `json:"download_limit" bson:"download_limit"`
	DownloadCount int64              `json:"-" bson:"download_count"`
	ExpiresAt     *time.Time         `json:"expires_at" bson:"expires_at"`
	CreatedAt     time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt     *time.Time         `json:"updated_at" bson:"updated_at"`
}
