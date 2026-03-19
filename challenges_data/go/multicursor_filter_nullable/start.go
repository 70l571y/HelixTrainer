//go:build ignore

package main

type UserPatch struct {
	ID       int
	Email    string // nullable
	Nickname string // nullable
	Status   string
	City     string // nullable
	Country  string
	Phone    string // nullable
	Score    int
}
