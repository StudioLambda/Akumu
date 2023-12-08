package main

import (
	"github.com/studiolambda/akumu"
	"github.com/studiolambda/akumu/http"
)

func handler(request http.Request) (response http.Response) {
	return response.Status(http.StatusCreated).JSON(map[string]string{
		"message": "Hello, World!",
	})
}

func main() {
	akumu := akumu.New()

	akumu.Get("/", handler)
	akumu.Start()
}
