package main

import (
	"time"

	"github.com/RadiusNetworks/go-beacon"
	"github.com/RadiusNetworks/go-beacon/advertiser"
)

func main() {
	urlBeacon, _ := beacon.NewEddystoneURLBeacon("https://radiusnetworks.com", -42)
	eddystoneURLParser := beacon.NewParser("eddystone_url", beacon.DefaultLayouts["eddystone_url"])
	advert := eddystoneURLParser.GenerateAd(urlBeacon)
	adv, _ := advertiser.New()
	for {
		adv.AdvertiseMfgData(0xfeaa, advert)
		time.Sleep(5 * time.Second)
	}
}
