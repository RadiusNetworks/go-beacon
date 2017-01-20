package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/RadiusNetworks/go-beacon"
	"github.com/RadiusNetworks/go-beacon/advertiser"
)

func main() {
	urlBeacon, _ := beacon.NewEddystoneURLBeacon("https://www.radiusnetworks.com", -42)
	eddystoneURLParser := beacon.NewParser("eddystone_url", beacon.DefaultLayouts["eddystone_url"])
	advert := eddystoneURLParser.GenerateAd(urlBeacon)
	adv, _ := advertiser.New()
	adv.AdvertiseServiceData(0xfeaa, advert)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT)
	signal.Notify(sigChan, syscall.SIGTERM)
	<-sigChan // wait for signal
}
