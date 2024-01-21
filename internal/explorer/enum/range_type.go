package enum

type ContentRange string

const (
	FullContentRange    ContentRange = "bytes (\\d+)-(\\d+)/(\\d+)"
	WrappedContentRange ContentRange = "bytes=(\\d+)-(\\d*)"
)
