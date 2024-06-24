package main

import (
	"context"
	"flag"
	"log"

	"github.com/milosgajdos/go-vocode"
)

var (
	phoneNr string
)

func init() {
	flag.StringVar(&phoneNr, "phone-nr", "", "phone number")
}

func main() {
	flag.Parse()

	if phoneNr == "" {
		log.Fatal("must specify phone number")
	}

	client := vocode.NewClient()
	ctx := context.Background()

	extActionReq := &vocode.CreateActionReq{
		ActionReq: vocode.ActionReq{
			Type: vocode.ActionExternal,
			Config: vocode.ExternalActionConfig{
				ProcessingMode: vocode.MutedProcessing,
				Name:           "Baseconfig",
				Description:    "Some description",
				URL:            "https://foobar.com",
				InputSchema:    map[string]any{},
			},
			Trigger: &vocode.FnCallTrigger{
				Type:   vocode.FnCallTriggerType,
				Config: map[string]any{},
			},
		},
	}

	res, err := client.CreateAction(ctx, extActionReq)
	if err != nil {
		log.Fatalf("failed creating external action: %v", err)
	}
	log.Printf("%#v", res)

	trCallActionReq := &vocode.CreateActionReq{
		ActionReq: vocode.ActionReq{
			Type: vocode.ActionTransferCall,
			Config: vocode.TransferCallActionConfig{
				PhoneNr: phoneNr,
			},
			Trigger: &vocode.FnCallTrigger{
				Type:   vocode.FnCallTriggerType,
				Config: map[string]any{},
			},
		},
	}

	res, err = client.CreateAction(ctx, trCallActionReq)
	if err != nil {
		log.Fatalf("failed creating external action: %v", err)
	}

	a, err := client.GetAction(ctx, res.ID)
	if err != nil {
		log.Fatalf("failed getting action %s: %v", res.ID, err)
	}
	log.Printf("got action: %v", a.ID)

	actions, err := client.ListActions(ctx, nil)
	if err != nil {
		log.Fatalf("failed listing actions: %v", err)
	}
	log.Printf("got %d actions: %#v", len(actions.Items), actions)
}
