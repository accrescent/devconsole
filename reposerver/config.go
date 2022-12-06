package main

type config struct {
	APIKey     string `toml:"api_key"`
	PublishDir string `toml:"publish_dir"`
}
