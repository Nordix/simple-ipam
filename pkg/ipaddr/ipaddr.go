// Package ipaddr provides basic iterator functions for ip-addresses
//
// Some code is taken from github.com/mikioh/ipaddr. The internal ipv6
// representation is exposed to provide a hashable (immutable) address
// type.
package ipaddr

import (
	"fmt"
	"math"
	"net"
)

// String stringifies an IPv6Int
func (i IPv6Int) String() string {
	return fmt.Sprintf("%016x,%016x", i[0], i[1])
}
func first(i IPv6Int, m net.IPMask) IPv6Int {
	mask := ipMaskToIPv6Int(m)
	return IPv6Int{i[0] & mask[0], i[1] & mask[1]}
}
func last(i IPv6Int, m net.IPMask) IPv6Int {
	mask := ipMaskToIPv6Int(m)
	return IPv6Int{i[0] | ^mask[0], i[1] | ^mask[1]}
}

// Cidr represents a CIDR range
type Cidr struct {
	Current, First, Last IPv6Int
	Size                 uint64
}

func to16(n *net.IPNet) *net.IPNet {
	o, b := n.Mask.Size()
	if b >= 128 {
		return n
	}
	return &net.IPNet{
		IP:   n.IP.To16(),
		Mask: net.CIDRMask(o+(128-b), 128),
	}
}

func rangeSize(n *net.IPNet) uint64 {
	o, _ := n.Mask.Size()
	if o <= 64 {
		return math.MaxUint64
	}
	return uint64(1) << uint(128-o)
}

// IPToIPv6Int convertes a net.IP to a IPv6Int
func IPToIPv6Int(ip net.IP) IPv6Int {
	return ipToIPv6Int(ip.To16())
}

// NewCidr creates a new Cidr. Error if an invalid cidr string is passed
func NewCidr(cidr string) (*Cidr, error) {
	_, net, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}
	net = to16(net)
	i := ipToIPv6Int(net.IP)
	current := first(i, net.Mask)
	return &Cidr{
		Current: current,
		First:   current,
		Last:    last(i, net.Mask),
		Size:    rangeSize(net),
	}, nil
}

// Step steps the Current address in the Cidr. Wraps to First after Last.
func (c *Cidr) Step() {
	if (&c.Current).Cmp(&c.Last) < 0 {
		(&c.Current).incr()
	} else {
		c.Current = c.First
	}
}

// IP returns the net.IP representation of a IPv6Int.
func IP(i IPv6Int) net.IP {
	ip := i.ip()
	if i4 := ip.To4(); i4 != nil {
		return i4
	}
	return ip
}
