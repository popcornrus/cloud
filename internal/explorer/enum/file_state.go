package enum

type FileState string

const (
	FileStatePending    FileState = "pending"
	FileStateCollecting FileState = "collecting"
	FileStateUploading  FileState = "uploading"
	FileStateConverting FileState = "converting"
	FileStateDone       FileState = "done"
	FileStateFailed     FileState = "failed"
)
