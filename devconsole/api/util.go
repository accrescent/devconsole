package api

import (
	"archive/zip"
	"bytes"
	"errors"
	"io"
	"mime/multipart"

	"github.com/accrescent/apkstat"
	"google.golang.org/protobuf/proto"

	pb "github.com/accrescent/devconsole/pb"
)

var ErrFatalIO = errors.New("fatal IO error")

func openAPKSet(
	formFile *multipart.FileHeader,
) (*pb.BuildApksResult, *apk.APK, multipart.File, error) {
	file, err := formFile.Open()
	if err != nil {
		return nil, nil, nil, err
	}

	apkSet, err := zip.NewReader(file, formFile.Size)
	if err != nil {
		return nil, nil, nil, err
	}

	rawBaseAPK, err := apkSet.Open("splits/base-master.apk")
	if err != nil {
		return nil, nil, nil, err
	}
	baseAPK, err := io.ReadAll(rawBaseAPK)
	if err != nil {
		return nil, nil, nil, ErrFatalIO
	}
	apk, err := apk.FromReader(bytes.NewReader(baseAPK), int64(len(baseAPK)))
	if err != nil {
		return nil, nil, nil, err
	}

	metadataFile, err := apkSet.Open("toc.pb")
	if err != nil {
		return nil, nil, nil, err
	}
	defer metadataFile.Close()
	metadataData, err := io.ReadAll(metadataFile)
	if err != nil {
		return nil, nil, nil, err
	}
	metadata := new(pb.BuildApksResult)
	if err := proto.Unmarshal(metadataData, metadata); err != nil {
		return nil, nil, nil, err
	}

	return metadata, apk, file, nil
}
