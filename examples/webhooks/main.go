package main

import (
	"context"
	"log"

	"github.com/milosgajdos/go-vocode"
)

func main() {
	client := vocode.NewClient()
	ctx := context.Background()

	whCreateReq := &vocode.CreateWebhookReq{
		WebhookReq: vocode.WebhookReq{
			Subs: []vocode.Event{
				vocode.MessageEvent,
				vocode.ActionEvent,
			},
			URL:    "https://foobar.com",
			Method: vocode.PostWebhook,
		},
	}

	res, err := client.CreateWebhook(ctx, whCreateReq)
	if err != nil {
		log.Fatalf("failed creating webhook: %v", err)
	}
	log.Printf("created webhook: %v", res)

	a, err := client.GetWebhook(ctx, res.ID)
	if err != nil {
		log.Fatalf("failed getting webhook %s: %v", res.ID, err)
	}
	log.Printf("got webhook: %v", a.ID)

	webhooks, err := client.ListWebhooks(ctx, nil)
	if err != nil {
		log.Fatalf("failed listing webhooks: %v", err)
	}
	log.Printf("got %d webhooks: %#v", len(webhooks.Items), webhooks)
}
