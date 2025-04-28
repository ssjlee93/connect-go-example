package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	greetv1 "example/gen/greet/v1"
	"example/gen/greet/v1/greetv1connect"

	"connectrpc.com/connect"
)

func main() {
	ctx := context.Background()
	client := greetv1connect.NewGreetServiceClient(
		http.DefaultClient,
		"http://localhost:8080",
		connect.WithGRPC(),
	)
	res, err := client.Greet(
		ctx,
		connect.NewRequest(&greetv1.GreetRequest{Name: "Unary user"}),
	)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(res.Msg.Greeting)
	// Client streaming
	clientStream := client.ClientGreet(ctx)
	clientStream.Send(&greetv1.GreetRequest{Name: "client stream"})
	clientStream.Send(&greetv1.GreetRequest{Name: "client stream user"})
	res, err = clientStream.CloseAndReceive()
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(res.Msg.Greeting)

	// Server streaming
	serverStream, err := client.ServerGreet(ctx, connect.NewRequest(&greetv1.GreetRequest{Name: "server stream user"}))
	if err != nil {
		log.Println(err)
		return
	}
	for serverStream.Receive() {
		log.Println(serverStream.Msg().Greeting)
	}

	// Bidi streaming
	stream := client.BothGreet(ctx)

	if err != nil {
		log.Fatalf("failed to open stream: %v", err)
	}
	defer stream.CloseRequest()

	namesToSend := []string{"Alice", "Bob", "Charlie", "David"}

	// Send a bunch of messages to the server
	log.Println("Sending messages to the server...")
	for _, name := range namesToSend {
		fmt.Println(name)
		req := &greetv1.GreetRequest{Name: name}
		if err := stream.Send(req); err != nil {
			log.Fatalf("failed to send request: %v", err)
		}
		log.Printf("Sent: %s", name)
		// time.Sleep(time.Millisecond * 200) // Simulate some delay between sends
	}

	// Signal to the server that we're done sending
	if err := stream.CloseRequest(); err != nil {
		log.Fatalf("failed to close request stream: %v", err)
	}
	log.Println("Finished sending messages.")

	// Expecting one greeting message back from the server
	log.Println("Waiting for greeting message from the server...")
	resp, err := stream.Receive()
	if err != nil {
		log.Fatalf("failed to receive response: %v", err)
	}
	log.Printf("Received greeting: %s", resp.Greeting)
}
