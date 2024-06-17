package main

import (
	"context"
	"flag"
	"log"

	"github.com/milosgajdos/go-vocode"
)

var (
	areaCode     string
	telProvider  string
	telAccountID string
)

func init() {
	flag.StringVar(&areaCode, "area-code", "951", "area code")
	flag.StringVar(&telProvider, "tel-provider", "twilio", "telephony provider [twilio or vonage]")
	flag.StringVar(&telAccountID, "tel-account-id", "", "telephony account ID")
}

func main() {
	flag.Parse()

	if telAccountID == "" {
		log.Fatal("telephony account ID can not be empty")
	}

	client := vocode.NewClient()
	ctx := context.Background()

	buyReq := &vocode.BuyNumberReq{
		AreaCode:     areaCode,
		TelProvider:  vocode.TelProvider(telProvider),
		TelAccountID: telAccountID,
	}

	number, err := client.BuyNumber(ctx, buyReq)
	if err != nil {
		log.Fatalf("failed buying a new number: %v", err)
	}

	updateReq := &vocode.UpdateNumberReq{
		Label: "Foobar",
	}
	if _, err := client.UpdateNumber(ctx, number.Number, updateReq); err != nil {
		log.Fatalf("failed updating number %s: %v", number.Number, err)
	}

	numbers, err := client.ListNumbers(ctx, nil)
	if err != nil {
		log.Fatalf("failed getting numbers: %v", err)
	}
	log.Printf("got numbers: %d", len(numbers.Items))

	if len(numbers.Items) > 0 {
		number, err := client.GetNumber(ctx, numbers.Items[0].Number)
		if err != nil {
			log.Fatalf("failed getting number: %v", err)
		}
		log.Printf("got number: %+v", number)
	}

	if _, err := client.CancelNumber(ctx, number.Number); err != nil {
		log.Fatalf("failed cancelling number %s: %v", number.Number, err)
	}
}
