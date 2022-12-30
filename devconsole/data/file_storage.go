package data

import (
	"io"
	"mime/multipart"
)

type FileStorage interface {
	SaveNewApp(
		apkSet multipart.File,
		icon multipart.File,
	) (apkSetHandle string, iconHandle string, err error)
	SaveUpdate(apkSet multipart.File) (apkSetHandle string, err error)

	GetAPKSet(apkSetHandle string) (file io.Reader, size int64, err error)

	// DeleteApp takes a file handle to an app and deletes it from disk
	DeleteApp(handle string) error
	// DeleteIcon takes a file handle to an icon and deletes it from disk
	DeleteIcon(handle string) error
}
