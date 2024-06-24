package main

import (
	"context"
	"log"

	"github.com/milosgajdos/go-vocode"
)

func main() {
	client := vocode.NewClient()
	ctx := context.Background()

	voiceAzureReq := &vocode.CreateVoiceReq{
		VoiceReq: vocode.VoiceReq{
			Type: vocode.AzureVoiceType,
			AzureVoice: &vocode.AzureVoice{
				Name:  "FoobarAzure",
				Pitch: 123,
				Rate:  123,
			},
		},
	}

	res, err := client.CreateVoice(ctx, voiceAzureReq)
	if err != nil {
		log.Fatalf("failed creating voice: %v", err)
	}
	log.Printf("created voice: %v", res)

	voiceRimeReq := &vocode.CreateVoiceReq{
		VoiceReq: vocode.VoiceReq{
			Type: vocode.RimeVoiceType,
			RimeVoice: &vocode.RimeVoice{
				Speaker:    "Frank",
				SpeedAlpha: 12,
				ModelID:    vocode.MistRimeVoiceModel,
			},
		},
	}

	res, err = client.CreateVoice(ctx, voiceRimeReq)
	if err != nil {
		log.Fatalf("failed creating voice: %v", err)
	}
	log.Printf("created voice: %v", res)

	a, err := client.GetVoice(ctx, res.ID)
	if err != nil {
		log.Fatalf("failed getting voice %s: %v", res.ID, err)
	}
	log.Printf("got voice: %v", a.ID)

	voices, err := client.ListVoices(ctx, nil)
	if err != nil {
		log.Fatalf("failed listing voices: %v", err)
	}
	log.Printf("got %d voices: %#v", len(voices.Items), voices)
}
