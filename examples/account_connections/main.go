package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/milosgajdos/go-vocode"
)

var (
	twillioAccountID string
)

func init() {
	flag.StringVar(&twillioAccountID, "twillio-account-id", "", "Twillio account ID")
}

func main() {
	client := vocode.NewClient()
	ctx := context.Background()

	oaiAPIKey := os.Getenv("OPENAI_API_KEY")
	if oaiAPIKey == "" {
		log.Fatal("missing openai API key")
	}

	twillioAuthToken := os.Getenv("TWILLIO_AUTH_TOKEN")
	if twillioAuthToken == "" {
		log.Fatal("missing twillio auth token")
	}

	oaiConnReq := &vocode.CreateAccountConnReq{
		AccountConnReq: vocode.AccountConnReq{
			Type: vocode.AccountConnOpenAI,
			OpenAIAccount: &vocode.OpenAIAccount{
				Creds: &vocode.OpenAICreds{
					APIKey: oaiAPIKey,
				},
			},
		},
	}

	res, err := client.CreateAccountConn(ctx, oaiConnReq)
	if err != nil {
		log.Fatalf("failed creating openai account connection: %v", err)
	}
	log.Printf("created openai account connection: %v", res)

	twillioConnReq := &vocode.CreateAccountConnReq{
		AccountConnReq: vocode.AccountConnReq{
			Type: vocode.AccountConnTwilio,
			TwilioAccount: &vocode.TwilioAccount{
				Creds: &vocode.TwilioCreds{
					AccountID: twillioAccountID,
					AuthToken: twillioAuthToken,
				},
			},
		},
	}

	res, err = client.CreateAccountConn(ctx, twillioConnReq)
	if err != nil {
		log.Fatalf("failed creating twilio account connection: %v", err)
	}
	log.Printf("created twillio account connection: %v", res)

	a, err := client.GetAccountConn(ctx, res.ID)
	if err != nil {
		log.Fatalf("failed getting account connection %s: %v", res.ID, err)
	}
	log.Printf("got account connection: %v", a.ID)

	accountConns, err := client.ListAccountConns(ctx, nil)
	if err != nil {
		log.Fatalf("failed listing account connections: %v", err)
	}
	log.Printf("got %d account connections: %#v", len(accountConns.Items), accountConns)
}
