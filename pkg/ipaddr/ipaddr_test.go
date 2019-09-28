package ipaddr

import (
	"net"
	"testing"
)

func create(t *testing.T, cidr, current, last string, size uint64) *Cidr {
	c, err := NevCidr(cidr)
	if err != nil {
		t.Errorf("Unexpected error")
	}
	if !IP(c.Current).Equal(net.ParseIP(current)) {
		t.Errorf("Current=%s, expected %s", IP(c.Current), current)
	}
	if c.Current.Cmp(&c.First) != 0 {
		t.Errorf("Current=%s, First=%s", IP(c.Current), IP(c.First))
	}
	t.Logf("Cidr=%s, First=%s, Last=%s\n", cidr, c.First, c.Last)
	if c.Size != size {
		t.Errorf("Size=%d, expected %d", c.Size, size)
	}
	return c
}

func step(t *testing.T, c *Cidr, current string) {
	c.Step()
	if !IP(c.Current).Equal(net.ParseIP(current)) {
		t.Errorf("Current=%s, expected %s", IP(c.Current), current)
	}
}

func TestBasic(t *testing.T) {
	c := create(t, "10.0.0.0/24", "10.0.0.0", "10.0.0.255", 256)
	step(t, c, "10.0.0.1")
	step(t, c, "10.0.0.2")
	step(t, c, "10.0.0.3")

	c = create(t, "1000::2222/128", "1000::2222", "1000::2222", 1)
	step(t, c, "1000::2222")
	step(t, c, "1000::2222")
	c = create(t, "10.0.0.22/32", "10.0.0.22", "10.0.0.22", 1)
	step(t, c, "10.0.0.22")
	step(t, c, "10.0.0.22")

	c = create(t, "10.10.10.8/30", "10.10.10.8", "10.10.10.11", 4)
	step(t, c, "10.10.10.9")
	step(t, c, "10.10.10.10")
	step(t, c, "10.10.10.11")
	step(t, c, "10.10.10.8")
	step(t, c, "10.10.10.9")

	c = create(t, "1000::4448/126", "1000::4448", "1000::444b", 4)
	step(t, c, "1000::4449")
	step(t, c, "1000::444a")
	step(t, c, "1000::444b")
	step(t, c, "1000::4448")
	step(t, c, "1000::4449")
}
