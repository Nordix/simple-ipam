package ipaddr

// The code in this file is taken from github.com/mikioh/ipaddr.
// The only change is that the type and the Cmp function has been made public.

// Copyright 2013 Mikio Hara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.

import (
	"encoding/binary"
	"math"
	"net"
)

// IPv6Int represents a ipv6 address. This type can be used as key in a map.
type IPv6Int [2]uint64

func (i *IPv6Int) incr() {
	if i[1] == math.MaxUint64 {
		i[0]++
		i[1] = 0
	} else {
		i[1]++
	}
}

func ipToIPv6Int(ip net.IP) IPv6Int {
	return IPv6Int{binary.BigEndian.Uint64(ip[:8]), binary.BigEndian.Uint64(ip[8:16])}
}

func (i *IPv6Int) ip() net.IP {
	ip := make(net.IP, net.IPv6len)
	binary.BigEndian.PutUint64(ip[:8], i[0])
	binary.BigEndian.PutUint64(ip[8:16], i[1])
	return ip
}

func ipMaskToIPv6Int(m net.IPMask) IPv6Int {
	return IPv6Int{binary.BigEndian.Uint64(m[:8]), binary.BigEndian.Uint64(m[8:16])}
}

// Cmp compares two IPv6Int
func (i *IPv6Int) Cmp(j *IPv6Int) int {
	if i[0] < j[0] {
		return -1
	}
	if i[0] > j[0] {
		return +1
	}
	if i[1] < j[1] {
		return -1
	}
	if i[1] > j[1] {
		return +1
	}
	return 0
}
