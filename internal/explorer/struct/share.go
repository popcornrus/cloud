package _struct

import (
	"time"
)

type (
	PinResult struct {
		Uuid      string    `json:"uuid"`
		PinCode   bool      `json:"pin_code"`
		CreatedAt time.Time `json:"created_at"`
	}

	ShareFileResult struct {
		Uuid    string  `json:"uuid"`
		Name    string  `json:"name"`
		Size    int64   `json:"size"`
		Type    string  `json:"type"`
		Preview *string `json:"preview"`
	}

	ShareResult struct {
		Uuid string          `json:"uuid"`
		File ShareFileResult `json:"file"`
	}
)
