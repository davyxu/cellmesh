package memsd

import "time"

type Config struct {
	Address        string
	RequestTimeout time.Duration
}

func DefaultConfig() Config {

	return Config{
		Address:        ":8900",
		RequestTimeout: time.Second * 10,
	}
}
