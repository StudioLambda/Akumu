package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/studiolambda/akumu"
	"github.com/studiolambda/akumu/middleware"
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

func fails(request *http.Request) error {
	return errors.New("something went wrong")
}

func fails2(request *http.Request) error {
	return akumu.
		Response(http.StatusNotFound).
		Failed(errors.New("something went wrong"))
}

func fails3(request *http.Request) error {
	return akumu.Failed(akumu.Problem{
		Type:     "http://example.com/problems/not-found",
		Title:    http.StatusText(http.StatusNotFound),
		Detail:   "The requested resource could not be found.",
		Status:   http.StatusNotFound,
		Instance: request.URL.String(),
	})
}

func fails4(request *http.Request) error {
	return akumu.Response(http.StatusUnauthorized).
		Failed(akumu.Problem{
			Type:     "http://example.com/problems/not-found",
			Title:    http.StatusText(http.StatusNotFound),
			Detail:   "The requested resource could not be found.",
			Status:   http.StatusNotFound,
			Instance: request.URL.String(),
		})
}

func fails5(request *http.Request) error {
	return akumu.Response(http.StatusNotFound).
		Failed(akumu.Problem{
			Type:     "http://example.com/problems/not-found",
			Title:    "page not found",
			Detail:   "the requested resource could not be found.",
			Instance: request.URL.String(),
		})
}

func panicable(request *http.Request) error {
	panic("something went wrong")
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

	router.Handle("GET /", akumu.Handler(user))
	router.Handle("GET /sse", akumu.Handler(sse))
	router.Handle("GET /fails", akumu.Handler(fails))
	router.Handle("GET /fails2", akumu.Handler(fails2))
	router.Handle("GET /fails3", akumu.Handler(fails3))
	router.Handle("GET /fails4", akumu.Handler(fails4))
	router.Handle("GET /fails5", akumu.Handler(fails5))
	router.Handle("GET /panic", akumu.Handler(panicable))

	server := http.Server{
		Addr:    ":8080",
		Handler: middleware.Recover(router),
	}

	fmt.Println("Server is running on: ", server.Addr)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}
