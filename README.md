# simple-ipam

A super-simple IPAM.

This IPAM administers addresses from a single CIDR range, e.g `1100::/120`.

Example;

```go
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
