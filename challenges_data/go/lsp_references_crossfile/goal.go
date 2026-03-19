//go:build ignore

package main

func buildToken(requestID string) string {
	return "token:" + requestID
}

func main() {
	requestID := "abc"
	println(buildToken(requestID))
	logAudit(requestID)
}
