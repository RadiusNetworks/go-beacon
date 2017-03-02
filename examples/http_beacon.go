package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/RadiusNetworks/go-beacon"
	"github.com/RadiusNetworks/go-beacon/advertiser"
)

type Beacon struct {
	BeaconType  string
	Identifiers struct {
		UUID  string
		Major uint16
		Minor uint16
	}
}

func main() {
	/*
		urlBeacon, _ := beacon.NewEddystoneURLBeacon("https://www.radiusnetworks.com", -42)
		eddystoneURLParser := beacon.NewParser("eddystone_url", beacon.DefaultLayouts["eddystone_url"])
		advert := eddystoneURLParser.GenerateAd(urlBeacon)
		adv, _ := advertiser.New()
		adv.AdvertiseServiceData(0xfeaa, advert)
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT)
		signal.Notify(sigChan, syscall.SIGTERM)
		<-sigChan // wait for signal
	*/
	http.HandleFunc("/beacon", beaconHandler)
	log.Fatal(http.ListenAndServe(":9999", nil))
}

func beaconHandler(response http.ResponseWriter, request *http.Request) {
	decoder := json.NewDecoder(request.Body)
	var beacon Beacon
	err := decoder.Decode(&beacon)
	if err != nil {
		panic(err)
	}
	defer request.Body.Close()
	advertiseBeacon(beacon)
}

func advertiseBeacon(advBeacon Beacon) {
	altBeaconParser := beacon.NewParser("altbeacon", beacon.DefaultLayouts["altbeacon"])
	altBeacon := beacon.NewAltBeacon(advBeacon.Identifiers.UUID, advBeacon.Identifiers.Major, advBeacon.Identifiers.Minor, -42)
	advert := altBeaconParser.GenerateAd(altBeacon)
	adv, _ := advertiser.New()
	adv.AdvertiseMfgData(0xbeef, advert)
	log.Println(fmt.Sprintf("Advertising beacon: UUID: %s, Major %d, Minor %d", advBeacon.Identifiers.UUID, advBeacon.Identifiers.Major, advBeacon.Identifiers.Minor))
}