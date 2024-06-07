package main

import (
	"context"
	"log"

	"github.com/milosgajdos/go-vocode"
)

func main() {
	client := vocode.NewClient()
	ctx := context.Background()

	extActionReq := &vocode.CreateReq{
		Type: vocode.External,
		Config: vocode.ExternalActionConfig{
			ProcessingMode: vocode.MutedProcessing,
			Name:           "Baseconfig",
			Description:    "Some description",
			URL:            "https://foobar.com",
			InputSchema:    map[string]any{},
		},
		Trigger: vocode.FnCallTrigger{
			Type:   vocode.FnCallTriggerType,
			Config: map[string]any{},
		},
	}

	res, err := client.CreateAction(ctx, extActionReq)
	if err != nil {
		log.Fatalf("failed creating external action: %v", err)
	}
	log.Printf("%#v", res)

	trCallActionReq := &vocode.CreateReq{
		Type: vocode.TransferCall,
		Config: vocode.TransferCallActionConfig{
			PhoneNr: "+19517449404",
		},
		Trigger: vocode.FnCallTrigger{
			Type:   vocode.FnCallTriggerType,
			Config: map[string]any{},
		},
	}

	res, err = client.CreateAction(ctx, trCallActionReq)
	if err != nil {
		log.Fatalf("failed creating external action: %v", err)
	}
	log.Printf("%#v", res)

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
