package data

import (
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

const APP_PATH = "app.apks"

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
	appDir, err := os.MkdirTemp(s.baseDir, "")
	if err != nil {
		return "", "", err
	}
	iconDir, err := os.MkdirTemp(s.baseDir, "")
	if err != nil {
		return "", "", err
	}

	appPath := filepath.Join(appDir, APP_PATH)
	if err := saveFile(apkSet, appPath); err != nil {
		return "", "", err
	}
	iconPath := filepath.Join(iconDir, "icon.png")
	if err := saveFile(icon, iconPath); err != nil {
		return "", "", err
	}

	return appPath, iconPath, nil
}

func saveFile(src multipart.File, dst string) error {
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)

	return err
}

func (s *LocalStorage) SaveUpdate(apkSet multipart.File) (apkSetHandle string, err error) {
	dir, err := os.MkdirTemp(s.baseDir, "")
	if err != nil {
		return "", err
	}

	appPath := filepath.Join(dir, APP_PATH)
	if err := saveFile(apkSet, appPath); err != nil {
		return "", err
	}

	return appPath, nil
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

func (s *LocalStorage) DeleteApp(handle string) error {
	return os.Remove(handle)
}

func (s *LocalStorage) DeleteIcon(handle string) error {
	return os.Remove(handle)
}
