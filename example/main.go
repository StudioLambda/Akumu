package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/studiolambda/akumu"
)

type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func user(request *http.Request) error {
	return akumu.
		Response(http.StatusOK).
		JSON(User{Name: "Erik C. Forés", Email: "soc@erik.cat"})
}

func sse(request *http.Request) error {
	messages := make(chan []byte)

	go func() {
		defer close(messages)

		for i := 0; i < 5; i++ {
			messages <- []byte(fmt.Sprintf("data: %d\n\n", i))
			time.Sleep(1 * time.Second)
		}
	}()

	return akumu.
		Response(http.StatusOK).
		SSE(messages)
}

func main() {
	router := http.NewServeMux()

	router.HandleFunc("GET /", akumu.Handle(user))
	router.HandleFunc("GET /sse", akumu.Handle(sse))

	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	fmt.Println("Server is running on: ", server.Addr)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}
