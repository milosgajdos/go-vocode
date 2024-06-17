package main

import (
	"context"
	"log"

	"github.com/milosgajdos/go-vocode"
)

func main() {
	client := vocode.NewClient()
	ctx := context.Background()

	extActionReq := &vocode.CreateActionReq{
		ActionReqBase: vocode.ActionReqBase{
			Type: vocode.ExternalActionType,
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
		},
	}

	res, err := client.CreateAction(ctx, extActionReq)
	if err != nil {
		log.Fatalf("failed creating external action: %v", err)
	}
	log.Printf("%#v", res)

	trCallActionReq := &vocode.CreateActionReq{
		ActionReqBase: vocode.ActionReqBase{
			Type: vocode.TransferCallActionType,
			Config: vocode.TransferCallActionConfig{
				PhoneNr: "+19517449404",
			},
			Trigger: vocode.FnCallTrigger{
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
