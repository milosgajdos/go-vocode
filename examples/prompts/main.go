package main

import (
	"context"
	"log"

	"github.com/milosgajdos/go-vocode"
)

func main() {
	client := vocode.NewClient()
	ctx := context.Background()

	createPromptReq := &vocode.CreatePromptReq{
		PromptReqBase: vocode.PromptReqBase{
			Content: "You are a voice assistant that answers questions about coding",
			Fields: []vocode.Field{
				{
					Type:  vocode.EmailFieldType,
					Label: "Foo",
					Name:  "FooPrompt",
					Desc:  "Example prompt",
				},
			},
		},
	}

	res, err := client.CreatePrompt(ctx, createPromptReq)
	if err != nil {
		log.Fatalf("failed creating prompt: %v", err)
	}
	log.Printf("created prompt: %v", res)

	a, err := client.GetPrompt(ctx, res.ID)
	if err != nil {
		log.Fatalf("failed getting prompt %s: %v", res.ID, err)
	}
	log.Printf("got prompt: %v", a.ID)

	prompts, err := client.ListPrompts(ctx, nil)
	if err != nil {
		log.Fatalf("failed listing prompts: %v", err)
	}
	log.Printf("got %d prompts: %#v", len(prompts.Items), prompts)
}
