//go:build ignore

package main

const lookupToken = "legacy_token"

func main() {
	println(loadToken(lookupToken))
}
