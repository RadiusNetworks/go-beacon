# Beacon package for Go

### Example: Scanning for beacons
> This package currently only works with a BLE112 to scan

```go
package main

import (
  "fmt"
  "github.com/RadiusNetworks/go-beacon"
  "github.com/RadiusNetworks/go-beacon/ble112"
  "time"
)


func main() {
  device := ble112.Devices()[0]
  scanner := beacon.NewScanner(device, beacon.DefaultParsers())
  data := make(chan beacon.Slice)
  done := make(chan bool)
  go scanner.Scan(time.Second, data, done)

  for {
    beacons := <- data
    for _, beacon := range beacons {
      fmt.Printf("%v\n", beacon)
    }
    fmt.Printf("\n")
  }
}
```

### http_beacon

`http_beacon` listens for an HTTP POST and broadcasts as the specified beacon.
Out of the box, the port is `9999`.  The format of the JSON (so far) is:
```json
{
  "beacon_type": "altbeacon",
  "identifiers": {
    "major": 1,
    "minor": 2,
    "uuid": "BCDB5AEB-F4E9-4600-B90F-70A2BE2F88C3"
  }
}
```

If you were to put the above into a `altbeacon.json` file, you could send this to a running `http_beacon` using the following `curl` command:

```
curl -X 'POST' -H 'Content-Type:application/json' -d @altbeacon.json 'http://<hostname>:9999/beacon'
```

In order to cross-compile `http_beacon` to run on a pi, use the following build commands:

```
sudo env GOOS=linux GOARCH=arm GOPATH=<path to your Go workspace> go get -v github.com/radiusnetworks/go-beacon/http_beacon
env GOOS=linux GOARCH=arm go build -v github.com/radiusnetworks/go-beacon/http_beacon
```

The `go get` command gets the necessary dependencies for the specified OS and architecture.  The `go build` command builds the executable and deposits it in the current working directory.  You can then copy this file onto your pi and run it using sudo:

```
sudo ./http_beacon
```
