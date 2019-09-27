# simple-ipam

A super-simple IPAM.

This IPAM administers addresses from a single CIDR range, e.g `1100::/120`.

[![GoDoc](https://godoc.org/github.com/Nordix/simple-ipam/pkg/ipam?status.svg)](https://godoc.org/github.com/Nordix/simple-ipam/pkg/ipam)

Example;

```go
package main

import (
	"fmt"
	"github.com/Nordix/simple-ipam/pkg/ipam"
)

func main() {
	cidr := "1100::/120"
	ipam, _ := ipam.New(cidr)
	fmt.Printf("Unallocated addresses in %s; %d\n", cidr, ipam.Unallocated())
	a, _ := ipam.Allocate()
	fmt.Printf("Allocated; %s\n", a)
	ipam.Free(a)
}
```

[Go playground](https://play.golang.org/p/2JAl0s9S5s9)
