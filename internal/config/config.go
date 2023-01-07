package config

import "time"

type Config struct {
	Host                  string
	Key                   string
	DataBaseAddress       string
	AccrualServiceAddress string
	ClientDuration        time.Duration
}

func New(serverAddress, key, dbAddress, accrualServiceAddress string) Config {
	return Config{
		Host:                  serverAddress,
		Key:                   key,
		DataBaseAddress:       dbAddress,
		AccrualServiceAddress: accrualServiceAddress,
		ClientDuration:        2 * time.Second,
	}
}
