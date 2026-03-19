//go:build ignore

package main

type bannerFormatter struct{}

func (bannerFormatter) Format(input string) string {
	return "[" + input + "]"
}
