package config

import "time"

type WatchConfig struct {
	EndPoints   string
	DialTimeout time.Duration
	WatchKey    string
	PostAddress string
}

type CalculateConfig struct {
	EndPoints      string
	DialTimeout    time.Duration
	UpdateInterval time.Duration
	ChannelCount   int
	DataKey        string
	DataValue      string
}
