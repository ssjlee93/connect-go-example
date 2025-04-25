package main

import (
	"context"
	"log"
	"net/http"

	greetv1 "example/gen/greet/v1"
	"example/gen/greet/v1/greetv1connect"
	"example/pkg/interceptor"

	"connectrpc.com/connect"
)

func main() {
	interceptors := connect.WithInterceptors(interceptor.NewAuthInterceptor())
	ctx := context.Background()
	client := greetv1connect.NewGreetServiceClient(
		http.DefaultClient,
		"http://localhost:8080",
		connect.WithGRPC(),
		interceptors,
	)
	res, err := client.Greet(
		ctx,
		connect.NewRequest(&greetv1.GreetRequest{Name: "Jane"}),
	)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(res.Msg.Greeting)
	// Client streaming
	stream := client.GreetClient(ctx)
	stream.Send(&greetv1.GreetRequest{Name: "John"})
	stream.Send(&greetv1.GreetRequest{Name: "Doe"})
	res, err = stream.CloseAndReceive()
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(res.Msg.Greeting)
}
