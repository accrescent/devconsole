package quality

type UploadType int

const (
	NewApp UploadType = iota
	AppUpdate
)

const (
	MIN_TARGET_SDK_NEW_APP    = 31
	MIN_TARGET_SDK_APP_UPDATE = 30
)
