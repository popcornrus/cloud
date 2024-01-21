package _struct

type (
	PrepareResult struct {
		Url       string `json:"url"`
		ChunkSize int64  `json:"chunk_size"`
	}

	PreviewProcessing struct {
		Width  string `json:"width"`
		Height string `json:"height"`
		Action string `json:"action"`
	}
)
