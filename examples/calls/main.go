package main

import (
	"context"
	"flag"
	"log"

	"github.com/milosgajdos/go-vocode"
)

func main() {
	flag.Parse()

	client := vocode.NewClient()
	ctx := context.Background()

	calls, err := client.ListCalls(ctx, nil)
	if err != nil {
		log.Fatalf("failed listing calls: %v", err)
	}
	log.Printf("got %d calls: %#v", len(calls.Items), calls)

	if len(calls.Items) > 0 {
		res, err := client.GetCall(ctx, calls.Items[0].ID)
		if err != nil {
			log.Fatalf("failed getting call %s: %v", calls.Items[0].ID, err)
		}
		log.Printf("call: %+v", res)
	}
}
