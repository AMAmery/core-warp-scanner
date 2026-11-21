package generator

import (
	"net/netip"
	"testing"
)

func TestParseCIDR_Valid(t *testing.T) {
	// IPv4
	v4, err := ParseCIDR("162.159.192.0/24")
	if err != nil {
		t.Fatalf("unexpected error parsing IPv4: %v", err)
	}
	if !v4.IsV4() {
		t.Error("expected IsV4() to be true for IPv4")
	}
	if v4.Mask() != 24 {
		t.Errorf("expected Mask() = 24, got %d", v4.Mask())
	}
	if v4.String() != "162.159.192.0/24" {
		t.Errorf("expected String() = 162.159.192.0/24, got %s", v4.String())
	}

	// IPv6
	v6, err := ParseCIDR("2606:4700:d1::/48")
	if err != nil {
		t.Fatalf("unexpected error parsing IPv6: %v", err)
	}
	if v6.IsV4() {
		t.Error("expected IsV4() to be false for IPv6")
	}
	if v6.Mask() != 48 {
		t.Errorf("expected Mask() = 48, got %d", v6.Mask())
	}
	if v6.String() != "2606:4700:d1::/48" {
		t.Errorf("expected String() = 2606:4700:d1::/48, got %s", v6.String())
	}

	// MustParseCIDR
	mustV4 := MustParseCIDR("1.1.1.0/24")
	if mustV4.String() != "1.1.1.0/24" {
		t.Errorf("unexpected result from MustParseCIDR: %s", mustV4.String())
	}
}

func TestParseCIDR_Invalid(t *testing.T) {
	invalid := []string{
		"",
		"not-a-cidr",
		"192.168.1.1",
		"192.168.1.1/33",
		"999.999.999.999/24",
	}

	for _, raw := range invalid {
		_, err := ParseCIDR(raw)
		if err == nil {
			t.Errorf("expected error for invalid CIDR %q, got nil", raw)
		}
	}
}

func TestCIDR_SplitTo24(t *testing.T) {
	// /24 -> 1 subnet
	c24 := MustParseCIDR("162.159.192.0/24")
	subs24, err := c24.SplitTo24()
	if err != nil {
		t.Fatalf("unexpected error splitting /24: %v", err)
	}
	if len(subs24) != 1 {
		t.Fatalf("expected 1 subnet, got %d", len(subs24))
	}
	if subs24[0].String() != "162.159.192.0/24" {
		t.Errorf("expected 162.159.192.0/24, got %s", subs24[0].String())
	}

	// /20 -> 16 subnets
	c20 := MustParseCIDR("188.114.96.0/20")
	subs20, err := c20.SplitTo24()
	if err != nil {
		t.Fatalf("unexpected error splitting /20: %v", err)
	}
	if len(subs20) != 16 {
		t.Fatalf("expected 16 subnets, got %d", len(subs20))
	}
	if subs20[0].String() != "188.114.96.0/24" {
		t.Errorf("expected first subnet 188.114.96.0/24, got %s", subs20[0].String())
	}
	if subs20[15].String() != "188.114.111.0/24" {
		t.Errorf("expected last subnet 188.114.111.0/24, got %s", subs20[15].String())
	}

	// /25 -> error (mask > 24)
	c25 := MustParseCIDR("192.168.1.0/25")
	if _, err := c25.SplitTo24(); err == nil {
		t.Error("expected error splitting /25 into /24, got nil")
	}

	// IPv6 -> error
	c6 := MustParseCIDR("2606:4700:d1::/48")
	if _, err := c6.SplitTo24(); err == nil {
		t.Error("expected error splitting IPv6 into /24, got nil")
	}
}

func TestNewCIDRPool(t *testing.T) {
	parent := MustParseCIDR("188.114.96.0/22")
	pool, err := NewCIDRPool(parent)
	if err != nil {
		t.Fatalf("unexpected error creating CIDRPool: %v", err)
	}

	if pool.Parent().String() != "188.114.96.0/22" {
		t.Errorf("expected Parent to be 188.114.96.0/22, got %s", pool.Parent().String())
	}

	// /22 gives 4 /24 subnets
	if len(pool.Subnets()) != 4 {
		t.Fatalf("expected 4 subnets in pool, got %d", len(pool.Subnets()))
	}

	// Stub test
	if ips := pool.Extract(10); ips != nil {
		t.Errorf("expected Extract stub to return nil, got %v", ips)
	}

	// Error when parent is IPv6
	v6Parent := MustParseCIDR("2606:4700:d1::/48")
	if _, err := NewCIDRPool(v6Parent); err == nil {
		t.Error("expected error creating CIDRPool with IPv6 parent, got nil")
	}
}

func TestSubnetBlock(t *testing.T) {
	prefix := netip.MustParsePrefix("10.0.0.0/24")
	block := NewSubnetBlock(prefix)

	if block.String() != "10.0.0.0/24" {
		t.Errorf("expected 10.0.0.0/24, got %s", block.String())
	}

	// Stub test
	if ips := block.ExtractUniqueIPs(5); ips != nil {
		t.Errorf("expected ExtractUniqueIPs stub to return nil, got %v", ips)
	}
}

func TestHostIP(t *testing.T) {
	// IPv4
	ip4, err := ParseHostIP("1.1.1.1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ip4.String() != "1.1.1.1" {
		t.Errorf("expected 1.1.1.1, got %s", ip4.String())
	}
	if ip4.WithPort(2408) != "1.1.1.1:2408" {
		t.Errorf("expected 1.1.1.1:2408, got %s", ip4.WithPort(2408))
	}

	// IPv6
	ip6, err := ParseHostIP("2606:4700::1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ip6.String() != "2606:4700::1" {
		t.Errorf("expected 2606:4700::1, got %s", ip6.String())
	}
	if ip6.WithPort(2408) != "[2606:4700::1]:2408" {
		t.Errorf("expected [2606:4700::1]:2408, got %s", ip6.WithPort(2408))
	}
}

