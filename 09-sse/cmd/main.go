package main

import (
	"context"
	"flag"
	"log"

	"github.com/kinesis-network/docker-image-samples/09-sse/sse"
)

var (
	mode *string = flag.String("mode", "s", "[s]erver or [c]lient")
	addr *string = flag.String("addr", "0.0.0.0:9000", "http endpoint")
)

func main() {
	flag.Parse()
	ctx := context.Background()

	if *mode == "s" {
		log.Println("Starting a server", *addr)
		sse.ServerMain(ctx, *addr)
		return
	}

	if *mode != "c" {
		log.Fatal("Unknown mode")
	}

	log.Println("Starting a client", *addr)
	sse.SubscribeToSse(ctx, *addr, func(s string) bool {
		log.Println(s)
		return true
	})
}
