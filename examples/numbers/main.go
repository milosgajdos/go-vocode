package main

import (
	"context"
	"log"

	"github.com/milosgajdos/go-vocode"
)

func main() {
	client := vocode.NewClient()

	numbers, err := client.ListNumbers(context.Background(), nil)
	if err != nil {
		log.Fatalf("failed getting numbers: %v", err)
	}

	log.Printf("got numbers: %d", len(numbers.Items))

	if len(numbers.Items) > 0 {
		number, err := client.GetNumber(context.Background(), numbers.Items[0].Number)
		if err != nil {
			log.Fatalf("failed getting number: %v", err)
		}
		log.Printf("got number: %v", number)
	}
}
