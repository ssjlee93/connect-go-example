package main

import (
	"example/gen/greet/v1/greetv1connect"
	"example/internal/service"
	"example/pkg/db"
	"net/http"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func main() {
	d := db.NewDB()
	// interceptor := connect.WithInterceptors(interceptor.NewAuthInterceptor())
	greeter := service.NewGreetServer(d)
	mux := http.NewServeMux()
	path, handler := greetv1connect.NewGreetServiceHandler(greeter)
	mux.Handle(path, handler)
	http.ListenAndServe(
		"localhost:8080",
		// Use h2c so we can serve HTTP/2 without TLS.
		h2c.NewHandler(mux, &http2.Server{}),
	)
}
