package main

import (
	"context"
	"fmt"
	"time"

	"github.com/studiolambda/akumu"
	"github.com/studiolambda/akumu/http"
	"github.com/studiolambda/golidate"
	"github.com/studiolambda/golidate/rule"
)

type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (user User) Validate(ctx context.Context) golidate.Results {
	return golidate.Validate(
		ctx,
		golidate.Value(user.Name).Name("name").Rules(
			rule.MinLen(2),
			rule.MaxLen(256),
		),
		golidate.Value(user.Email).Name("email").Rules(
			rule.Email(),
		),
	)
}

func show(request http.Request) (response http.Response) {
	return response.JSON(User{
		Name:  "John Doe",
		Email: "foo@example.com",
	})
}

func create(request http.Request) (response http.Response) {
	var user User

	if err := request.Validate(&user); err != nil {
		return response.Error(err)
	}

	// ...

	return response.Status(http.StatusCreated).JSON(user)
}

func sse(request http.Request) (response http.Response) {
	return response.HTML(`
	<!DOCTYPE html>
	<html lang="en">
	  <head>
		<meta charset="UTF-8" />
		<meta http-equiv="X-UA-Compatible" content="IE=edge" />
		<meta name="viewport" content="width=device-width, initial-scale=1.0" />
		<title>SSE</title>
	  </head>
	  <body onload="onLoaded()">
		<h1>Server Sent Events</h1>
		<div id="main">
		  <h3 id="rand-container">Random Number</h3>
		  <div id="random-number"></div>
		</div>
		<script>
		  const onLoaded = () => {
			let eventSource = new EventSource("/stream");
			eventSource.onmessage = (event) => {
			  document.getElementById("random-number").innerHTML = event.data;
			};
		  };
		</script>
	  </body>
	</html>
	`)
}

func stream(request http.Request) (response http.Response) {
	messages := make(chan http.SSE, 5)

	go func() {
		defer close(messages)

		for i := 0; i < 5; i++ {
			messages <- http.SSE{
				Data: []byte(fmt.Sprintf("%d", i)),
			}
			time.Sleep(1 * time.Second)
		}
	}()

	return response.SSE(messages)
}

func failure(request http.Request) (response http.Response) {
	err := http.NewError(nil).
		Status(http.StatusNotFound).
		Header("X-Error", "true")

	return response.Status(http.StatusNotFound).Error(err)
}

func main() {
	akumu := akumu.New()

	akumu.Get("/", show)
	akumu.Post("/", create)
	akumu.Get("/failure", failure)
	akumu.Get("/sse", sse)
	akumu.Get("/stream", stream)
	akumu.Start()
}
