//go:build ignore

package main

type Runner interface {
	Execute(input string) string
}
