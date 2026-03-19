//go:build ignore

package main

func buildToken(sessionID string) string {
	return "token:" + sessionID
}

func main() {
	sessionID := "abc"
	println(buildToken(sessionID))
	logAudit(sessionID)
}
