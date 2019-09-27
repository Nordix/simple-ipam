// Copyright 2019 Nordix foundation

// Package ipam is a very simple IPAM
//
// Administers a single CIDR range, e.g "1000::/124".
//
// The functions are NOT thread safe.
package ipam


import (
	"fmt"
	"math"
	"math/big"
	"net"

	"github.com/mikioh/ipaddr"
)

// IPAM holds the ipam state
type IPAM struct {
	// The original CIDR range
	CIDR net.IPNet
	prefix *ipaddr.Prefix
	cursor *ipaddr.Cursor
	allocated map[string]bool
}

// New creates a new IPAM for the passed CIDR.
// Error if the passed CIDR is invalid.
func New(cidr string) (*IPAM, error) {
	_, net, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}
	prefix := ipaddr.NewPrefix(net)
	cursor := ipaddr.NewCursor([]ipaddr.Prefix{*prefix})
	return &IPAM{
		CIDR:   *net,
		prefix: prefix,
		cursor: cursor,
		allocated: make(map[string]bool),
	}, nil
}

// Allocate allocates a new address.
// An error is returned if there is no addresses left.
func (i *IPAM) Allocate() (net.IP, error) {
	if i.Unallocated() < 1 {
		return nil, fmt.Errorf("No addresses left")
	}
	for {
		p := i.cursor.Pos()
		if i.cursor.Next() == nil {
			i.cursor.Reset(nil)
		}
		if _, ok := i.allocated[p.IP.String()]; !ok {
			i.allocated[p.IP.String()] = true
			return p.IP, nil
		}
	}
}

// Free frees an allocated address.
// To free a non-allocated address is a no-op.
func (i *IPAM) Free(a net.IP) {
	delete(i.allocated, a.String())
}

// Unallocated returns the number of unallocated addresses.
// If the number is > math.MaxInt64 then math.MaxInt64 is returned.
func (i *IPAM) Unallocated() int64 {
	tot := i.prefix.NumNodes()
	free := tot.Sub(tot, big.NewInt(int64(len(i.allocated))))
	if free.IsInt64() {
		return free.Int64()
	}
	return math.MaxInt64
}

// Reserve reserves an address.
// Error if the address is outside the CIDR or if the address is allocated already.
func (i *IPAM) Reserve(a net.IP) error {
	if ! i.CIDR.Contains(a) {
		return fmt.Errorf("Address outside the cidr")
	}
	if _, ok := i.allocated[a.String()]; ok {
		return fmt.Errorf("Address already allocated")
	}
	i.allocated[a.String()] = true
	return nil
}

// ReserveFirstAndLast reserves the first and last address.
// These are valid addresses but some programs may refuse to use them.
// Note that the number of Unallocated addresses may become zero.
func (i *IPAM) ReserveFirstAndLast() {
	i.Reserve(i.cursor.First().IP)
	i.Reserve(i.cursor.Last().IP)
}
