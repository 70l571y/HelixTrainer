//go:build ignore

package main

func buildProfile(email string) string {
	return normalizeEmail(email)
}

func main() {
	println(buildProfile("  User@Example.COM  "))
}
