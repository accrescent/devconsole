//go:build !debug

package main

import (
	"os"

	"github.com/BurntSushi/toml"
)

func loadConfig(path string) (*config, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var conf config
	err = toml.Unmarshal(file, &conf)

	return &conf, err
}
