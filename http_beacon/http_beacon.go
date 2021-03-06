package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/RadiusNetworks/go-beacon"
	"github.com/RadiusNetworks/go-beacon/advertiser"
)

// BeaconSpecification contains beacon type and configuration for advertising.
// Struct and fields must be exported for the json Decoder to work.
type BeaconSpecification struct {
	BeaconType  string
	Identifiers struct {
		UUID  string
		Major uint16
		Minor uint16
	}
}

func main() {
	http.HandleFunc("/beacon", beaconHandler)
	log.Fatal(http.ListenAndServe(":9999", nil))
}

func beaconHandler(response http.ResponseWriter, request *http.Request) {
	decoder := json.NewDecoder(request.Body)
	var beacon BeaconSpecification
	err := decoder.Decode(&beacon)
	if err != nil {
		panic(err)
	}
	defer request.Body.Close()
	advertiseBeacon(beacon)
}

func advertiseBeacon(advBeacon BeaconSpecification) {
	altBeaconParser := beacon.NewParser("altbeacon", beacon.DefaultLayouts["altbeacon"])
	altBeacon := beacon.NewAltBeacon(advBeacon.Identifiers.UUID, advBeacon.Identifiers.Major, advBeacon.Identifiers.Minor, -42)
	advert := altBeaconParser.GenerateAd(altBeacon)
	adv, _ := advertiser.New()
	adv.AdvertiseMfgData(0xbeef, advert)
	log.Println(fmt.Sprintf("Advertising beacon: UUID: %s, Major %d, Minor %d", advBeacon.Identifiers.UUID, advBeacon.Identifiers.Major, advBeacon.Identifiers.Minor))
}
