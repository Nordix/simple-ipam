
package ipam

import (
	"fmt"
	"math"
	"math/big"
	"net"

	"github.com/mikioh/ipaddr"
)

type IPAM struct {
	CIDR net.IPNet
	prefix *ipaddr.Prefix
	cursor *ipaddr.Cursor
	allocated map[string]bool
}

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

func (i *IPAM) Free(a net.IP) {
	delete(i.allocated, a.String())
}

func (i *IPAM) Unallocated() int64 {
	tot := i.prefix.NumNodes()
	free := tot.Sub(tot, big.NewInt(int64(len(i.allocated))))
	if free.IsInt64() {
		return free.Int64()
	}
	return math.MaxInt64
}
