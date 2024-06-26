package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/milosgajdos/go-vocode"
)

var (
	pcEnv   string
	pcIndex string
)

func init() {
	flag.StringVar(&pcEnv, "pinecone-env", "", "pinecone environment")
	flag.StringVar(&pcIndex, "pinecone-index", "", "pinecone index")
}

func main() {
	client := vocode.NewClient()
	ctx := context.Background()

	pcAPIKey := os.Getenv("PINECONE_API_KEY")
	if pcAPIKey == "" {
		log.Fatal("missing pinecone API key")
	}

	whCreateReq := &vocode.CreateVectorDBReq{
		VectorDBReq: vocode.VectorDBReq{
			Type:   vocode.PineConeVectorDB,
			Index:  pcIndex,
			APIKey: pcAPIKey,
			APIEnv: pcEnv,
		},
	}

	res, err := client.CreateVectorDB(ctx, whCreateReq)
	if err != nil {
		log.Fatalf("failed creating vectordb: %v", err)
	}
	log.Printf("created vectordb: %v", res)

	a, err := client.GetVectorDB(ctx, res.ID)
	if err != nil {
		log.Fatalf("failed getting vectordb %s: %v", res.ID, err)
	}
	log.Printf("got vectordb: %v", a.ID)

	vectordbs, err := client.ListVectorDBs(ctx, nil)
	if err != nil {
		log.Fatalf("failed listing vectordbs: %v", err)
	}
	log.Printf("got %d vectordbs: %#v", len(vectordbs.Items), vectordbs)
}
