package main

import (
	"context"
	"fmt"
	"github.com/disc/twirpexample/rpc/haberdasher"
	"github.com/twitchtv/twirp"
	"net/http"
	"os"
)

func main() {
	client := haberdasher.NewHaberdasherProtobufClient("http://localhost:8080", &http.Client{})

	hat, err := client.MakeHat(context.Background(), &haberdasher.Size{Inches: 5})
	if err != nil {
		twerr := err.(twirp.Error)
		fmt.Printf("oh no: %v\n", err)
		fmt.Println("retryable", twerr.Meta("retryable"))
		fmt.Println("retry_after", twerr.Meta("retry_after"))
		os.Exit(1)
	}
	fmt.Printf("I have a nice new hat: %+v", hat)
}
