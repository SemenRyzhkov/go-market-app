package config

type Config struct {
	Host            string
	Key             string
	DataBaseAddress string
}

func New(serverAddress, key, dbAddress string) Config {
	return Config{
		Host:            serverAddress,
		Key:             key,
		DataBaseAddress: dbAddress,
	}
}
