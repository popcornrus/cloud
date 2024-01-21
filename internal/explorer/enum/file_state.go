package enum

type FileState string

const (
	FileStatePending    FileState = "pending"
	FileStateCollecting FileState = "collecting"
	FileStateUploading  FileState = "uploading"
	FileStateDone       FileState = "done"
	FileStateFailed     FileState = "failed"
)
