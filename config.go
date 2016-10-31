package main

import (
	_ "log"
)

const (
	URL        = "https://github.com/Sangdol/pazuzu-test-repo.git"
	BASE_IMAGE = "ubuntu:14.04"
)

var config *Config

type GitConfig struct {
	Url string
}

type Config struct {
	ConfigType string
	Base       string
	Git        *GitConfig
}

func NewConfig() error {
	// TODO: add read from $HOME/.pazuzu/config and return error if fail
	// viper library is planned to be used here
	config = &Config{ConfigType: "git",
		Base: BASE_IMAGE,
		Git:  &GitConfig{Url: URL}}

	return nil
}
