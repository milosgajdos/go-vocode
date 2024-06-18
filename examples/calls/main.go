package main

import (
	"context"
	"flag"
	"log"

	"github.com/milosgajdos/go-vocode"
)

var (
	fromNr string
	toNr   string
	agent  string
)

func init() {
	flag.StringVar(&fromNr, "from-nr", "", "from number")
	flag.StringVar(&toNr, "to-nr", "", "to number")
	flag.StringVar(&agent, "agent", "", "agent ID")
}

func main() {
	flag.Parse()

	client := vocode.NewClient()
	ctx := context.Background()

	createCallReq := &vocode.CreateCallReq{
		FromNr:          fromNr,
		ToNr:            toNr,
		Agent:           agent,
		OnHumanNoAnswer: vocode.HangupCallOnNoHumanAnswer,
	}

	res, err := client.CreateCall(ctx, createCallReq)
	if err != nil {
		log.Fatalf("failed creating call: %v", err)
	}
	log.Printf("created call: %v", res)

	calls, err := client.ListCalls(ctx, nil)
	if err != nil {
		log.Fatalf("failed listing calls: %v", err)
	}
	log.Printf("got %d calls: %#v", len(calls.Items), calls)

	if len(calls.Items) > 0 {
		res, err := client.GetCall(ctx, calls.Items[0].ID)
		if err != nil {
			log.Fatalf("failed getting call %s: %v", calls.Items[0].ID, err)
		}
		log.Printf("call: %+v", res)
	}
}
