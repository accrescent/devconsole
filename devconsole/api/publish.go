package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/accrescent/devconsole/config"
	"github.com/accrescent/devconsole/data"
	"github.com/accrescent/devconsole/quality"
)

func publish(
	c *gin.Context, appID string, versionCode int32, versionName string,
	uploadType quality.UploadType, appFileHandle string,
) error {
	storage := c.MustGet("storage").(data.FileStorage)
	conf := c.MustGet("config").(config.Config)

	var method string
	if uploadType == quality.NewApp {
		method = http.MethodPost
	} else if uploadType == quality.Update {
		method = http.MethodPut
	}

	file, size, err := storage.GetAPKSet(appFileHandle)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return err
	}

	req, err := http.NewRequest(
		method, fmt.Sprintf("%s/apps/%s/%d/%s", conf.RepoURL, appID, versionCode, versionName),
		file,
	)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return err
	}
	req.Header.Add("Authorization", "token "+conf.APIKey)
	req.ContentLength = size
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		var err error
		switch resp.StatusCode {
		case http.StatusBadRequest:
			err = errors.New("bad request")
			c.AbortWithStatus(http.StatusInternalServerError)
		case http.StatusUnauthorized:
			err = errors.New("invalid repo server API key")
			_ = c.AbortWithError(http.StatusInternalServerError, err)
		case http.StatusConflict:
			err = errors.New("app already published")
			_ = c.AbortWithError(resp.StatusCode, err)
		default:
			err = errors.New("unknown error")
			_ = c.AbortWithError(http.StatusInternalServerError, err)
		}
		return err
	}

	return nil
}
