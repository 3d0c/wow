package config

type Config struct {
	Addr                  string `json:"addr"`          // e.g.: "127.0.0.1:5050"
	Timeout               int    `json:"timeout"`       // RW timeout in seconds
	MaxOpen               int    `json:"max_open"`      // Maximum opened connections
	HashcashDuration      int64  `json:"hash_duration"` // In seconds
	HashcashMaxIterations int    `json:"hash_max_iter"` // max iterations to prevent stuck on hard hashes (only for client)
	ZeroCount             int    `json:"zero_count"`    // Count of leading zeros
}
