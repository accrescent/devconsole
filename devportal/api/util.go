package api

import (
	"archive/zip"
	"bytes"
	"errors"
	"io"

	"github.com/accrescent/apkstat"
)

var ErrFatalIO = errors.New("fatal IO error")

func apkFromAPKSet(filename string) (*apk.APK, error) {
	apkSet, err := zip.OpenReader(filename)
	if err != nil {
		return nil, err
	}
	defer apkSet.Close()
	rawBaseAPK, err := apkSet.Open("splits/base-master.apk")
	if err != nil {
		return nil, err
	}
	baseAPK, err := io.ReadAll(rawBaseAPK)
	if err != nil {
		return nil, ErrFatalIO
	}

	apk, err := apk.FromReader(bytes.NewReader(baseAPK), int64(len(baseAPK)))
	if err != nil {
		return nil, err
	}

	return apk, nil
}
