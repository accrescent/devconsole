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

	// DeleteFile takes a file handle and deletes the associated file
	DeleteFile(handle string) error
}
