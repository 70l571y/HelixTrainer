package main

func main() {
	serverName := "api"
	println(serverName)
	println(buildURL(serverName))
	replicaCount := 3
	println(replicaCount)
}

func buildURL(name string) string {
	return "https://" + name
}
