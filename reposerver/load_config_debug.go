//go:build debug

package main

import "os"

func loadConfig(path string) (*config, error) {
	conf := &config{
		APIKey:     os.Getenv("API_KEY"),
		PublishDir: os.Getenv("PUBLISH_DIR"),
	}

	return conf, nil
}
