# go-vocode

[![Build Status](https://github.com/milosgajdos/go-vocode/workflows/CI/badge.svg)](https://github.com/milosgajdos/go-vocode/actions?query=workflow%3ACI)
[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/milosgajdos/go-vocode)
[![License: Apache-2.0](https://img.shields.io/badge/License-Apache--2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

A Go module for [vocode.dev](https://www.vocode.dev/) API client.

The official [vocode.dev](https://www.vocode.dev/)  API documentation, upon which this Go module has been built, can be found [here](https://docs.vocode.dev/api-reference).

# Get Started

Build the module

```shell
go build ./...
```

Run tests
```shell
go test ./...
```

You must sign up for the Vocode account and generate the API key, first, before you can use the API.

There are a few code samples available in the [examples](./examples) directory so please do have a look. They could give you an idea about how to use this Go module.

> [!IMPORTANT]
> Before you attempt to run the samples you must set an environment variable with the API key.
> These are automatically read by the client when it gets created; you can override them in your own code.

* `VOCODE_API_KEY`: Vocode API key


## Nix

There is a [Nix flake](https://nixos.wiki/wiki/Flakes) file available which lets you work on the Go module using nix.

Just run the following command and you are in the business:
```shell
nix develop
```

# Basics

Vocode lets you create conversational agents and make them available via a phone number.

Vocode currently provides two phone providers:
* [Twilio](https://www.twilio.com/en-us)
* [Vonage](https://www.vonage.co.uk/)

You must create an account with either before you can create a Vocode agent.

You can buy a phone number via the Vocode API, specifically the [Buy Number API endpoint](https://docs.vocode.dev/api-reference/numbers/buy-number).

You wil associate this number with an agent you'll create later on.

You can then proceed with configuring your agent:
* create a [voice config](./examples/voices)
* create conversational [prompt config](./examples/prompts)
* configure [agent actions](./examples/actions)

Once you've set up the phone number you can create an agent like so:

```Go
package main

import (
	"context"
	"flag"
	"log"
	"strings"

	"github.com/milosgajdos/go-vocode"
)

var (
	prompt  string
	voice   string
	actions string
)

func init() {
	flag.StringVar(&prompt, "prompt", "", "prompt ID")
	flag.StringVar(&voice, "voice", "", "voice ID")
	flag.StringVar(&actions, "actions", "", "comma separated list of actions")
}

func main() {
	flag.Parse()

	client := vocode.NewClient()
	ctx := context.Background()

	ax := strings.Split(actions, ",")

	createAgentReq := &vocode.CreateAgentReq{
		AgentReq: vocode.AgentReq{
			Name:                     "My Agent",
			Prompt:                   prompt,
			Voice:                    voice,
			Language:                 vocode.English,
			InterruptSense:           vocode.LowInterruptSense,
			EndpointSense:            vocode.AutoEndpointSense,
			IVRNavMode:               vocode.OffIVRMode,
			Speed:                    1.0,
			AsktIfHumanPresentOnIdle: true,
			Actions:                  ax,
		},
	}

	res, err := client.CreateAgent(ctx, createAgentReq)
	if err != nil {
		log.Fatalf("failed creating agent: %v", err)
	}
	log.Printf("created agent: %v", res)
}
```
