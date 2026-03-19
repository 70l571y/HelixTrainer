//go:build ignore

package main

type Config struct {
	APITimeout    int
	WorkerTimeout int
	APIRetries    int
	JobRetries    int
}

var cfg = Config{
	APITimeout:    60,
	WorkerTimeout: 60,
	APIRetries:    2,
	JobRetries:    3,
}
