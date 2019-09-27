
package ipam

import (
	"math"
	"net"
	"testing"
)

func allocate(t *testing.T, ipam *IPAM, expected string, left int64) {
	a, err := ipam.Allocate()
	if err != nil {
		if expected != "" {
			t.Errorf("Unexpected error for %s", expected)
		}
		return
	}
	if !a.Equal(net.ParseIP(expected)) {
		t.Errorf("Address %s, expected %s", a, expected)
	}
	u := ipam.Unallocated()
	if u != left {
		t.Errorf("Unallocated %d, expected %d", u, left)
	}
}
func free(t *testing.T, ipam *IPAM, addr string, left int64) {
	ipam.Free(net.ParseIP(addr))
	u := ipam.Unallocated()
	if u != left {
		t.Errorf("Unallocated %d, expected %d", u, left)
	}	
}
func create(t *testing.T, cidr string, left int64) *IPAM {
	ipam, err := New(cidr)
	if err != nil {
		t.Errorf("Failed to create ipam %s", cidr)
	}
	i := ipam.Unallocated()
	if i != left {
		t.Errorf("Unallocated %d, expected %d", i, left)
	}
	return ipam
}


func TestBasic(t *testing.T) {
	ipam, err := New("malformed")
	if err == nil {
		t.Errorf("Could create a malformed ipam")
	}

	ipam = create(t, "1000::/127", 2)
	allocate(t, ipam, "1000::", 1)
	allocate(t, ipam, "1000::1", 0)
	allocate(t, ipam, "", 0)
	free(t, ipam, "1000::", 1)
	free(t, ipam, "1000::", 1)
	free(t, ipam, "1000::", 1)
	allocate(t, ipam, "1000::", 0)

	ipam = create(t, "10.10.10.0/29", 8)
	allocate(t, ipam, "10.10.10.0", 7)
	allocate(t, ipam, "10.10.10.1", 6)
	allocate(t, ipam, "10.10.10.2", 5)
	allocate(t, ipam, "10.10.10.3", 4)
	allocate(t, ipam, "10.10.10.4", 3)
	allocate(t, ipam, "10.10.10.5", 2)
	allocate(t, ipam, "10.10.10.6", 1)
	allocate(t, ipam, "10.10.10.7", 0)

	free(t, ipam, "10.10.10.3", 1)
	free(t, ipam, "10.10.10.5", 2)

	allocate(t, ipam, "10.10.10.3", 1)

	free(t, ipam, "10.10.10.0", 2)
	free(t, ipam, "10.10.10.1", 3)
	free(t, ipam, "10.10.10.2", 4)
	free(t, ipam, "10.10.10.3", 5)
	allocate(t, ipam, "10.10.10.5", 4)

	ipam = create(t, "1000::1000/128", 1)
	ipam = create(t, "1000::/64", math.MaxInt64)
	ipam = create(t, "1000::/65", math.MaxInt64)
	allocate(t, ipam, "1000::", math.MaxInt64)
	allocate(t, ipam, "1000::1", math.MaxInt64-1)
	allocate(t, ipam, "1000::2", math.MaxInt64-2)
}

