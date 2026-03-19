//go:build ignore

package main

type Config struct {
	APITimeout    int
	WorkerTimeout int
	APIRetries    int
	JobRetries    int
}

var cfg = Config{
	APITimeout:    30,
	WorkerTimeout: 45,
	APIRetries:    2,
	JobRetries:    3,
}
