package data

import (
	"io"
	"mime/multipart"
	"os"
)

type LocalStorage struct {
	baseDir string
}

func NewLocalStorage(baseDir string) *LocalStorage {
	return &LocalStorage{baseDir}
}

func (s *LocalStorage) SaveNewApp(
	apkSet multipart.File,
	icon multipart.File,
) (apkSetHandle string, iconHandle string, err error) {
	appFile, err := os.CreateTemp(s.baseDir, "*.apks")
	if err != nil {
		return "", "", err
	}
	defer appFile.Close()
	iconFile, err := os.CreateTemp(s.baseDir, "*.png")
	if err != nil {
		return "", "", err
	}
	defer iconFile.Close()

	if _, err := io.Copy(appFile, apkSet); err != nil {
		return "", "", err
	}
	if _, err := io.Copy(iconFile, icon); err != nil {
		return "", "", err
	}

	return appFile.Name(), iconFile.Name(), nil
}

func (s *LocalStorage) SaveUpdate(apkSet multipart.File) (apkSetHandle string, err error) {
	appFile, err := os.CreateTemp(s.baseDir, "*.apks")
	if err != nil {
		return "", err
	}
	defer appFile.Close()

	if _, err := io.Copy(appFile, apkSet); err != nil {
		return "", err
	}

	return appFile.Name(), nil
}

func (s *LocalStorage) GetAPKSet(apkSetHandle string) (file io.Reader, size int64, err error) {
	apkSet, err := os.Open(apkSetHandle)
	if err != nil {
		return nil, 0, err
	}
	apkSetInfo, err := apkSet.Stat()
	if err != nil {
		return nil, 0, err
	}

	return apkSet, apkSetInfo.Size(), nil
}

func (s *LocalStorage) DeleteFile(handle string) error {
	return os.Remove(handle)
}
