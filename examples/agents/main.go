package main

import (
	"context"
	"log"

	"github.com/milosgajdos/go-vocode"
)

func main() {
	client := vocode.NewClient()
	ctx := context.Background()

	agents, err := client.ListAgents(ctx, nil)
	if err != nil {
		log.Fatalf("failed listing agents: %v", err)
	}
	log.Printf("got %d agents: %#v", len(agents.Items), agents)

	if len(agents.Items) > 0 {
		res, err := client.GetAgent(ctx, agents.Items[0].ID)
		if err != nil {
			log.Fatalf("failed getting agent %s: %v", agents.Items[0].ID, err)
		}
		log.Printf("agent: %+v", res)
	}
}
