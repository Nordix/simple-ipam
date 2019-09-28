// GO111MODULE=on go build -o /tmp/xtest ./examples/export.go
// /tmp/xtest | jq
package main

import (
	"os"
	"encoding/json"
	"github.com/Nordix/simple-ipam/pkg/ipam"
)

func main() {
	cidr := "12.0.0.0/29"
	ipam, _ := ipam.New(cidr)
	ipam.ReserveFirstAndLast()
	a, _ := ipam.Allocate()
	a, _ = ipam.Allocate()
	ipam.Free(a)
	a, _ = ipam.Allocate()
	a, _ = ipam.Allocate()
	ipam.Free(a)
	a, _ = ipam.Allocate()
	json.NewEncoder(os.Stdout).Encode(ipam.Export())
}
