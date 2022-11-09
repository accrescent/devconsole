package api

import (
	"archive/zip"
	"bytes"
	"errors"
	"io"
	"mime/multipart"
	"os"

	"github.com/accrescent/apkstat"
)

var ErrFatalIO = errors.New("fatal IO error")

func openAPKSet(formFile *multipart.FileHeader) (*apk.APK, multipart.File, error) {
	file, err := formFile.Open()
	if err != nil {
		return nil, nil, err
	}

	apkSet, err := zip.NewReader(file, formFile.Size)
	if err != nil {
		return nil, nil, err
	}
	rawBaseAPK, err := apkSet.Open("splits/base-master.apk")
	if err != nil {
		return nil, nil, err
	}
	baseAPK, err := io.ReadAll(rawBaseAPK)
	if err != nil {
		return nil, nil, ErrFatalIO
	}

	apk, err := apk.FromReader(bytes.NewReader(baseAPK), int64(len(baseAPK)))
	if err != nil {
		return nil, nil, err
	}

	return apk, file, nil
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
