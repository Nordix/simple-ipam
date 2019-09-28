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
	"net"

	"github.com/Nordix/simple-ipam/pkg/ipaddr"
)

// IPAM holds the ipam state
type IPAM struct {
	// The original CIDR range
	CIDR      net.IPNet
	cidr      *ipaddr.Cidr
	saveAddr  bool
	allocated map[ipaddr.IPv6Int]bool
}

// IPAMData data for export/import
type IPAMData struct {
	CIDR      string            `json:"cidr"`
	Current   string            `json:"current"`
	Flags     map[string]string `json:"flags,omitempty"`
	Allocated []string          `json:"allocated,omitempty"`
}

// Import re-creates a ipam from previously exported data
func Import(data *IPAMData) (*IPAM, error) {
	i, err := New(data.CIDR)
	if err != nil {
		return nil, err
	}
	i.cidr.Current = ipaddr.IPToIPv6Int(net.ParseIP(data.Current))
	for _, a := range data.Allocated {
		i.allocated[ipaddr.IPToIPv6Int(net.ParseIP(a))] = true
	}
	return i, nil
}

// Export export the ipam in a storable form
func (i *IPAM) Export() *IPAMData {
	var data IPAMData
	data.CIDR = i.CIDR.String()
	data.Current = ipaddr.IP(i.cidr.Current).String()
	data.Allocated = make([]string, len(i.allocated))
	x := 0
	for k := range i.allocated {
		data.Allocated[x] = ipaddr.IP(k).String()
		x = x + 1
	}
	return &data
}

// New creates a new IPAM for the passed CIDR.
// Error if the passed CIDR is invalid.
func New(cidr string) (*IPAM, error) {
	c, err := ipaddr.NewCidr(cidr)
	if err != nil {
		return nil, err
	}
	_, net, _ := net.ParseCIDR(cidr)
	return &IPAM{
		CIDR:      *net,
		cidr:      c,
		saveAddr:  c.Size < math.MaxUint64,
		allocated: make(map[ipaddr.IPv6Int]bool),
	}, nil
}

// Allocate allocates a new address.
// An error is returned if there is no addresses left.
func (i *IPAM) Allocate() (net.IP, error) {
	if !i.saveAddr {
		p := i.cidr.Current
		i.cidr.Step()
		return ipaddr.IP(p), nil
	}

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
	if !i.saveAddr {
		return
	}
	delete(i.allocated, ipaddr.IPToIPv6Int(a))
}

// Unallocated returns the number of unallocated addresses.
// If the number is > math.MaxUint64 then math.MaxUint64 is returned.
func (i *IPAM) Unallocated() uint64 {
	if !i.saveAddr {
		return i.cidr.Size
	}
	return i.cidr.Size - uint64(len(i.allocated))
}

// Reserve reserves an address.
// Error if the address is outside the CIDR or if the address is allocated already.
func (i *IPAM) Reserve(a net.IP) error {
	if !i.saveAddr {
		return fmt.Errorf("Not saving addresses so can't reserve")
	}
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
	if !i.saveAddr {
		if i.cidr.Current.Cmp(&i.cidr.First) == 0 {
			i.cidr.Step()
			return
		}
	}
	i.allocated[i.cidr.First] = true
	i.allocated[i.cidr.Last] = true
}
