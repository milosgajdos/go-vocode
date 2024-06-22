package main

import (
	"context"
	"log"

	"github.com/milosgajdos/go-vocode"
)

func main() {
	client := vocode.NewClient()

	usage, err := client.GetUsage(context.Background())
	if err != nil {
		log.Fatalf("failed getting usage: %v", err)
	}
	log.Printf("got usage: %+v", usage)
}
