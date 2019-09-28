// Copyright 2019 Nordix foundation

// Package ipam is a very simple IPAM
//
// Administers a single CIDR range, e.g "1000::/124".
//
// The functions are NOT thread safe.
package ipam

import (
	"fmt"
	"net"

	"github.com/Nordix/simple-ipam/pkg/ipaddr"
)

// IPAM holds the ipam state
type IPAM struct {
	// The original CIDR range
	CIDR      net.IPNet
	cidr      *ipaddr.Cidr
	allocated map[ipaddr.IPv6Int]bool
}

// New creates a new IPAM for the passed CIDR.
// Error if the passed CIDR is invalid.
func New(cidr string) (*IPAM, error) {
	c, err := ipaddr.NevCidr(cidr)
	if err != nil {
		return nil, err
	}
	_, net, _ := net.ParseCIDR(cidr)
	return &IPAM{
		CIDR:      *net,
		cidr:      c,
		allocated: make(map[ipaddr.IPv6Int]bool),
	}, nil
}

// Allocate allocates a new address.
// An error is returned if there is no addresses left.
func (i *IPAM) Allocate() (net.IP, error) {
	if i.Unallocated() < 1 {
		return nil, fmt.Errorf("No addresses left")
	}
	for {
		p := i.cidr.Current
		i.cidr.Step()
		if _, ok := i.allocated[p]; !ok {
			i.allocated[p] = true
			return ipaddr.IP(p), nil
		}
	}
}

// Free frees an allocated address.
// To free a non-allocated address is a no-op.
func (i *IPAM) Free(a net.IP) {
	delete(i.allocated, ipaddr.IPToIPv6Int(a))
}

// Unallocated returns the number of unallocated addresses.
// If the number is > math.MaxUint64 then math.MaxUint64 is returned.
func (i *IPAM) Unallocated() uint64 {
	if i.cidr.Size > 0 {
		return i.cidr.Size - uint64(len(i.allocated))
	}
	return 0
}

// Reserve reserves an address.
// Error if the address is outside the CIDR or if the address is allocated already.
func (i *IPAM) Reserve(a net.IP) error {
	if !i.CIDR.Contains(a) {
		return fmt.Errorf("Address outside the cidr")
	}
	ip := ipaddr.IPToIPv6Int(a)
	if _, ok := i.allocated[ip]; ok {
		return fmt.Errorf("Address already allocated")
	}
	i.allocated[ip] = true
	return nil
}

// ReserveFirstAndLast reserves the first and last address.
// These are valid addresses but some programs may refuse to use them.
// Note that the number of Unallocated addresses may become zero.
func (i *IPAM) ReserveFirstAndLast() {
	i.allocated[i.cidr.First] = true
	i.allocated[i.cidr.Last] = true
}
