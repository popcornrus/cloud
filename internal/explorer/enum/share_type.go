package enum

type ShareType int

const (
	BurnType ShareType = iota + 1
	DownloadsLimitType
	InfiniteType
)
