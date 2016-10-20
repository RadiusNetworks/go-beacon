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
  data := make(chan beacon.BeaconSlice)
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
