package main

import (
	"context"
	"log"
	"net/http"
	"sync"

	greetv1 "example/gen/greet/v1"
	"example/gen/greet/v1/greetv1connect"

	"connectrpc.com/connect"
)

var wg sync.WaitGroup

func main() {
	wg := sync.WaitGroup{}
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
	// TODO fix this
	stream := client.BothGreet(ctx)
	ch := make(chan bool)

	wg.Add(1)
	go send(ch, stream)

	<-ch

	go receive(ch, stream)

	wg.Wait()
}

func receive(state chan bool, stream *connect.BidiStreamForClient[greetv1.GreetRequest, greetv1.GreetResponse]) {
	defer wg.Done()
	// Expecting one greeting message back from the server
	log.Println("Waiting for greeting message from the server...")
	resp, err := stream.Receive()
	if err != nil {
		log.Fatalf("failed to receive response: %v", err)
	}
	log.Printf("Received greeting: %s", resp.Greeting)
	state <- false
}

func send(state chan bool, stream *connect.BidiStreamForClient[greetv1.GreetRequest, greetv1.GreetResponse]) {
	defer wg.Done()
	req := &greetv1.GreetRequest{Name: "bidi"}
	if err := stream.Send(req); err != nil {
		log.Fatalf("failed to send request: %v", err)
	}
	log.Printf("Sent: %s", "bidi")
	state <- true
}
