package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/digitalocean/godo"
	"github.com/kstkn/woody/goip"
)

var ip string

func main() {
	token := os.Getenv("DIGITALOCEAN_ACCESS_TOKEN")
	if token == "" {
		log.Fatal("Environment variable DIGITALOCEAN_ACCESS_TOKEN is not set (https://cloud.digitalocean.com/account/api/tokens)")
	}

	d, err := time.ParseDuration("WOODY_PERIOD")
	if err != nil {
		d = 5 * time.Minute
	}
	fmt.Printf("running with %s update interval", d)

	ticker := time.NewTicker(d)
	tickerDone := make(chan bool)
	go func() {
		for {
			select {
			case <-tickerDone:
				return
			case <-ticker.C:
				// find out current ip
				loc, err := goip.NewClient().GetLocation()
				if err != nil {
					log.Fatal(err)
				}
				fmt.Printf("current public IP: %s update interval", loc.Query)

				if loc.Query == ip {
					fmt.Print("public IP hasn't changed public IP")
					// do nothing, ip hasn't changed
				}

				// store and update ip
				ip = loc.Query

				do := godo.NewFromToken(os.Getenv("DIGITALOCEAN_ACCESS_TOKEN"))

				r, _, _ := do.Domains.RecordsByTypeAndName(context.Background(), "kstkn.com", "A", "pp.kstkn.com", nil)
				if len(r) != 1 {
					log.Fatal("no records found for update")
				}

				do.Domains.EditRecord(context.Background(), "kstkn.com", r[0].ID, &godo.DomainRecordEditRequest{
					Data: ip,
				})
				fmt.Print("DNS records updated")
			}
		}
	}()
}
