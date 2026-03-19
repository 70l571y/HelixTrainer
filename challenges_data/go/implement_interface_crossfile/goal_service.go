//go:build ignore

package main

type FileStore struct{}

func (FileStore) Save() string {
	return "saved"
}
