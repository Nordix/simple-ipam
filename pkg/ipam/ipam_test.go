package ipam

import (
	"math"
	"net"
	"testing"
)

func allocate(t *testing.T, ipam *IPAM, expected string, left uint64) {
	a, err := ipam.Allocate()
	u := ipam.Unallocated()
	if u != left {
		t.Fatalf("Unallocated %d, expected %d", u, left)
	}
	if err != nil {
		if expected != "" {
			t.Fatalf("Unexpected error for %s", expected)
		}
		return
	}
	if !a.Equal(net.ParseIP(expected)) {
		t.Fatalf("Address %s, expected %s", a, expected)
	}
}
func free(t *testing.T, ipam *IPAM, addr string, left uint64) {
	ipam.Free(net.ParseIP(addr))
	u := ipam.Unallocated()
	if u != left {
		t.Fatalf("Unallocated %d, expected %d", u, left)
	}
}
func create(t *testing.T, cidr string, left uint64) *IPAM {
	ipam, err := New(cidr)
	if err != nil {
		t.Fatalf("Failed to create ipam %s", cidr)
	}
	i := ipam.Unallocated()
	if i != left {
		t.Fatalf("Unallocated %d, expected %d", i, left)
	}
	cidrStr := ipam.CIDR.String()
	if cidrStr != cidr {
		t.Fatalf("CIDR set to %s, expected %s", cidrStr, cidr)
	}
	return ipam
}
func reserve(t *testing.T, ipam *IPAM, addr string, expectedErr bool, left uint64) {
	err := ipam.Reserve(net.ParseIP(addr))
	if err != nil {
		if !expectedErr {
			t.Fatalf("Unexpected error for %s", addr)
		}
	}
	u := ipam.Unallocated()
	if u != left {
		t.Fatalf("Unallocated %d, expected %d", u, left)
	}
}
func reserveFirstAndLast(t *testing.T, ipam *IPAM, left uint64) {
	ipam.ReserveFirstAndLast()
	u := ipam.Unallocated()
	if u != left {
		t.Fatalf("Unallocated %d, expected %d", u, left)
	}
}

func TestBasic(t *testing.T) {
	ipam, err := New("malformed")
	if err == nil {
		t.Fatalf("Could create a malformed ipam")
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
	ipam = create(t, "1000::/16", math.MaxUint64)
	ipam = create(t, "1000::/64", math.MaxUint64)
	ipam = create(t, "1000::/65", math.MaxInt64+1)
	allocate(t, ipam, "1000::", math.MaxInt64)
	allocate(t, ipam, "1000::1", math.MaxInt64-1)
	allocate(t, ipam, "1000::2", math.MaxInt64-2)

	ipam = create(t, "1000::/126", 4)
	reserve(t, ipam, "1000::", false, 3)
	reserve(t, ipam, "1000::", true, 3)
	reserve(t, ipam, "1000::3", false, 2)
	reserve(t, ipam, "1000::4", true, 2)

	ipam = create(t, "1000::/128", 1)
	reserveFirstAndLast(t, ipam, 0)
	ipam = create(t, "1000::/127", 2)
	reserveFirstAndLast(t, ipam, 0)
	ipam = create(t, "1000::/126", 4)
	reserveFirstAndLast(t, ipam, 2)
	allocate(t, ipam, "1000::1", 1)
	allocate(t, ipam, "1000::2", 0)
	free(t, ipam, "1000::2", 1)
	allocate(t, ipam, "1000::2", 0)

	ipam = create(t, "1000::/64", math.MaxUint64)
	reserveFirstAndLast(t, ipam, math.MaxUint64)
	allocate(t, ipam, "1000::1", math.MaxUint64)
	allocate(t, ipam, "1000::2", math.MaxUint64)
	allocate(t, ipam, "1000::3", math.MaxUint64)
	free(t, ipam, "1000::1", math.MaxUint64)
	reserve(t, ipam, "1000::4", true, math.MaxUint64)
}
