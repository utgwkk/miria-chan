package main

import (
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	err := postMessage("Hello from Golang")
	if err != nil {
		panic(err)
	}
}
