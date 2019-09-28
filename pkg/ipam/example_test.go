package ipam_test

import (
	"fmt"
	"github.com/Nordix/simple-ipam/pkg/ipam"
)

func Example() {
	cidr := "1100::/120"
	ipam, _ := ipam.New(cidr)
	fmt.Printf("Unallocated addresses in %s; %d\n", cidr, ipam.Unallocated())
	a, _ := ipam.Allocate()
	fmt.Printf("Allocated; %s\n", a)
	ipam.Free(a)

	// Output:
	// Unallocated addresses in 1100::/120; 256
	// Allocated; 1100::
}
