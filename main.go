package main

import (
	"context"
	"log"
	"os"

	"github.com/digitalocean/godo"
	"github.com/kstkn/woody/goip"
)

func main() {
	loc, err := goip.NewClient().GetLocation()
	if err != nil {
		log.Fatal(err)
	}

	do := godo.NewFromToken(os.Getenv("DIGITALOCEAN_ACCESS_TOKEN"))

	r, _, _ := do.Domains.RecordsByTypeAndName(context.Background(), "kstkn.com", "A", "pp.kstkn.com", nil)
	if len(r) != 1 {
		log.Fatal("no records found for update")
	}

	do.Domains.EditRecord(context.Background(), "kstkn.com", r[0].ID, &godo.DomainRecordEditRequest{
		Data: loc.Query,
	})
}
