package quality

type UploadType int

const (
	NewApp UploadType = iota
	Update
)

const (
	MIN_TARGET_SDK_NEW_APP = 33
	MIN_TARGET_SDK_UPDATE  = 33
)

const MIN_BUNDLETOOL_VERSION = "1.11.4"
