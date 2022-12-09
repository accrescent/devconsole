package quality

type UploadType int

const (
	NewApp UploadType = iota
	Update
)

const (
	MIN_TARGET_SDK_NEW_APP = 31
	MIN_TARGET_SDK_UPDATE  = 31
)

const MIN_BUNDLETOOL_VERSION = "1.13.1"
