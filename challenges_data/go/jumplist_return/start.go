//go:build ignore

package main

type Config struct {
	ServiceName string
	RetryCount  int
}

func loadConfig() Config {
	return Config{
		ServiceName: "api-gatway",
		RetryCount:  3,
	}
}

func warmup() {
	println("boot")
	println("cache")
	println("queue")
	println("http")
	println("workers")
	println("metrics")
	println("ready")
}

func main() {
	cfg := loadConfig()
	legasyTimeout := 30
	println(cfg.ServiceName, legasyTimeout)
	warmup()
}
