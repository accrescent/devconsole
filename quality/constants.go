package quality

type uploadType int

const (
	NewApp uploadType = iota
	AppUpdate
)

const (
	MIN_TARGET_SDK_NEW_APP    = 31
	MIN_TARGET_SDK_APP_UPDATE = 30
)
