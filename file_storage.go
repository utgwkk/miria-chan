package main

type FileStorage interface {
	Put(string) error
}
