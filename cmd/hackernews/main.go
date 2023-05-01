package main

import (
	"os"

	"github.com/captain-corgi/go-graphql-example/internal/hackernews"
)

func main() {
	os.Setenv("PORT", "8081")
	hackernews.Run()
}
