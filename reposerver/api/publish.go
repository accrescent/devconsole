package api

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func publish(c *gin.Context, uploadType uploadType) {
	publishDir := c.MustGet("publish_dir").(string)
	appID := c.Param("id")
	versionCode := c.Param("versionCode")
	versionCodeInt, err := strconv.Atoi(versionCode)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	version := c.Param("version")

	rawBody, err := io.ReadAll(c.Request.Body)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	body := bytes.NewReader(rawBody)
	apkSet, err := zip.NewReader(body, c.Request.ContentLength)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	repoData := repoData{
		Version:       version,
		VersionCode:   versionCodeInt,
		ABISplits:     nil,
		DensitySplits: nil,
		LangSplits:    nil,
	}

	// Extract APKs from APK set
	appDir := filepath.Join(publishDir, appID)
	apkOutDir := filepath.Join(appDir, versionCode)
	if err := os.MkdirAll(apkOutDir, 0755); err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	for _, file := range apkSet.File {
		path := filepath.Join(apkOutDir, file.Name)
		// zip.FileHeader.Name is not sanitized by default, so we have to sanitize it
		// ourselves.
		// See https://github.com/golang/go/issues/25849 and
		// https://github.com/uber/astro/pull/47.
		if !strings.HasPrefix(path, filepath.Clean(apkOutDir)+string(os.PathSeparator)) {
			// You're probably looking at the line below and asking yourself, "Why are
			// you using this status code?" The short answer is simple: I want to. The
			// long and more technical answer is that no other status code fits much
			// better. This is technically a client error since the developer pushed an
			// APK set with invalid/malicious ZIP file names, but the developer console
			// accepted that long ago. We still ought to sanitize the file names here as
			// opposed to strictly in the developer console because we don't need to
			// give an attacker with control over the developer console a vector toward
			// full control over the repository server. Labeling this as an internal
			// server error _technically might be correct but handling 500's in client
			// JS feels gross and wrong.
			//
			// So here we are. Do not extract files, do not brew coffee, do not pass go,
			// do not collect two hundred dollars.
			_ = c.AbortWithError(http.StatusTeapot, errors.New("path sanitization failed"))
			return
		}

		// Publish split APKs
		//
		// file is safe to read at this point
		if filepath.Ext(file.Name) == ".apk" {
			f, err := file.Open()
			if err != nil {
				_ = c.AbortWithError(http.StatusInternalServerError, err)
				return
			}
			defer f.Close()

			name, typ, typeName := getSplitInfo(filepath.Base(path))

			switch typ {
			case abi:
				repoData.ABISplits = append(repoData.ABISplits, typeName)
			case density:
				repoData.DensitySplits = append(repoData.DensitySplits, typeName)
			case lang:
				repoData.LangSplits = append(repoData.LangSplits, typeName)
			}

			newPath := filepath.Join(apkOutDir, name)
			split, err := os.OpenFile(newPath, os.O_WRONLY|os.O_CREATE, 0644)
			if err != nil {
				_ = c.AbortWithError(http.StatusInternalServerError, err)
				return
			}
			if _, err := io.Copy(split, f); err != nil {
				_ = c.AbortWithError(http.StatusInternalServerError, err)
				return
			}
		}
	}

	// Publish app repodata
	repoDataPath := filepath.Join(appDir, "repodata.json")
	openFlags := os.O_WRONLY
	if uploadType == newApp {
		openFlags |= os.O_CREATE
	} else if uploadType == appUpdate {
		openFlags |= os.O_TRUNC
	}
	repoDataFile, err := os.OpenFile(repoDataPath, openFlags, 0644)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	enc := json.NewEncoder(repoDataFile)
	if err := enc.Encode(repoData); err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if uploadType == appUpdate {
		// Delete old split APKs
		appDirPath := filepath.Join(publishDir, appID)
		appDir, err := os.Open(appDirPath)
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		subDirs, err := appDir.Readdirnames(-1)
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		for _, dir := range subDirs {
			if num, err := strconv.Atoi(dir); err == nil && num < versionCodeInt {
				// Directory presumably contains old split APKs. Delete them.
				if err := os.RemoveAll(filepath.Join(appDirPath, dir)); err != nil {
					_ = c.AbortWithError(http.StatusInternalServerError, err)
					return
				}
			}
		}
	}

	c.String(http.StatusOK, "")
}
