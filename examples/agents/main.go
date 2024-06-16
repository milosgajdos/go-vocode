package main

import (
	"context"
	"flag"
	"log"

	"github.com/milosgajdos/go-vocode"
)

var (
	prompt string
	voice  string
)

func init() {
	flag.StringVar(&prompt, "prompt", "", "prompt ID")
	flag.StringVar(&voice, "voice", "", "voice ID")
}

func main() {
	flag.Parse()

	client := vocode.NewClient()
	ctx := context.Background()

	createAgentReq := &vocode.CreateAgentReq{
		AgentReqbase: vocode.AgentReqbase{
			Name:                     "My Agent",
			Prompt:                   prompt,
			Voice:                    voice,
			Language:                 vocode.English,
			InterruptSense:           vocode.LowInterruptSense,
			EndpointSense:            vocode.AutoEndpointSense,
			IVRNavMode:               vocode.OffIVRMode,
			Speed:                    1.0,
			AsktIfHumanPresentOnIdle: true,
		},
	}

	res, err := client.CreateAgent(ctx, createAgentReq)
	if err != nil {
		log.Fatalf("failed creating agent: %v", err)
	}
	log.Printf("created agent: %v", res)

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
