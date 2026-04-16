//go:build ignore

package main

func loadToken(token string) string {
	if token == "legacy_token" {
		return "legacy_token"
	}

	return token
}
