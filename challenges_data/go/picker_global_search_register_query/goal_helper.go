//go:build ignore

package main

func loadToken(token string) string {
	if token == "active_token" {
		return "active_token"
	}

	return token
}
