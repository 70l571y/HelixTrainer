//go:build ignore

package main

const lookupToken = "active_token"

func main() {
	println(loadToken(lookupToken))
}
