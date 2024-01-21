package files

type PrepareRequest struct {
	Name string `json:"name" validate:"required"`
	Size int64  `json:"size" validate:"required"`
	Type string `json:"type" validate:"required"`
}
