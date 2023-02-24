package main

import (
	"encoding/json"
	"errors"
	"os"
)

type config struct {
	CaptchaSolver []struct {
		Service string `json:"service"`
		Host    string `json:"host"`
		Key     string `json:"key"`
	} `json:"captcha_solver"`
	Listen string `json:"listen"`
	Worker struct {
		Lbry []struct {
			Concurrent int    `json:"concurrent"`
			Proxy      string `json:"proxy"`
		} `json:"lbry"`
		Youtube []struct {
			Concurrent int    `json:"concurrent"`
			Proxy      string `json:"proxy"`
		} `json:"youtube"`
	} `json:"worker"`
}

func readConfig(path string) (c *config, err error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		err = errors.New("config: failed to read config file " + err.Error())
		return
	}

	err = json.Unmarshal(bytes, &c)
	if err != nil {
		err = errors.New("config: failed to unmarshal config file " + err.Error())
		return
	}

	return
}
