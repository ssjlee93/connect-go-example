package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
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

	wg.Add(2)
	go func() {
		defer wg.Done()
		send(stream)
		log.Println("finished sending")
	}()

	go func() {
		defer wg.Done()
		receive(stream)
		log.Println("finished receiving")
	}()

	wg.Wait()
	log.Println("client finished")
}

func receive(stream *connect.BidiStreamForClient[greetv1.GreetRequest, greetv1.GreetResponse]) {
	log.Println("receive message")
	for {
		// Expecting one greeting message back from the server
		log.Println("Waiting for greeting message from the server...")
		resp, err := stream.Receive()
		if err != nil {
			log.Printf("failed to receive response: %v", err)
		}
		if resp == nil {
			log.Println("received nil response")
			return // Exit the loop if response is nil
		}
		log.Printf("Received greeting: %s", resp.Greeting)
	}
}

func send(stream *connect.BidiStreamForClient[greetv1.GreetRequest, greetv1.GreetResponse]) {
	log.Println("send message")
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter your name:")

	for {
		input, _ := reader.ReadString('\n')
		input = input[:len(input)-1]
		if input == "exit()" {
			break
		}

		req := &greetv1.GreetRequest{Name: input}
		if err := stream.Send(req); err != nil {
			log.Printf("failed to send request: %v", err)
		}
		log.Printf("Sent: %s", input)
	}
}
