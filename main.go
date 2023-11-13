package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/digitalocean/godo"
	"github.com/kstkn/woody/goip"
)

var (
	token string
	ip    string
)

func update() {
	// find out current ip
	loc, err := goip.NewClient().GetLocation()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("current public IP: %s\n", loc.Query)

	if loc.Query == ip {
		log.Println("public IP hasn't changed public IP")
		// do nothing, ip hasn't changed
	}

	// store and update ip
	ip = loc.Query

	do := godo.NewFromToken(token)

	r, _, _ := do.Domains.RecordsByTypeAndName(context.Background(), "kstkn.com", "A", "pp.kstkn.com", nil)
	if len(r) != 1 {
		log.Fatal("no DNS records found for update")
	}

	if r[0].Data == ip {
		log.Println("DNS record already has correct IP")
		return
	}

	do.Domains.EditRecord(context.Background(), "kstkn.com", r[0].ID, &godo.DomainRecordEditRequest{
		Data: ip,
	})
	log.Println("DNS records updated")
}

func main() {
	token = os.Getenv("DIGITALOCEAN_ACCESS_TOKEN")
	if token == "" {
		log.Fatal("Environment variable DIGITALOCEAN_ACCESS_TOKEN is not set (https://cloud.digitalocean.com/account/api/tokens)")
	}

	d, err := time.ParseDuration("WOODY_PERIOD")
	if err != nil {
		d = 5 * time.Minute
	}
	log.Printf("running with %s update interval\n", d)

	signalNotifyContext, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	update()

	ticker := time.NewTicker(d)
	tickerDone := make(chan bool)
	go func() {
		for {
			select {
			case <-tickerDone:
				return
			case <-ticker.C:
				update()
			}
		}
	}()

	<-signalNotifyContext.Done()
}
