//go:build ignore

package main

type Config struct {
	ServiceName string
	RetryCount  int
}

func loadConfig() Config {
	return Config{
		ServiceName: "api-gateway",
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
	legacyTimeout := 30
	println(cfg.ServiceName, legacyTimeout)
	warmup()
}
