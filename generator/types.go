package generator

import (
	"encoding/binary"
	"fmt"
	"net/netip"
	"strings"
)

// CIDR represents a parent network range parsed from input files.
type CIDR struct {
	prefix netip.Prefix
}

// ParseCIDR parses and strictly validates a raw CIDR string using netip.ParsePrefix.
func ParseCIDR(raw string) (CIDR, error) {
	trimmed := strings.TrimSpace(raw)
	prefix, err := netip.ParsePrefix(trimmed)
	if err != nil {
		return CIDR{}, fmt.Errorf("invalid CIDR %q: %w", raw, err)
	}
	return CIDR{prefix: prefix.Masked()}, nil
}

// MustParseCIDR parses a raw CIDR string and panics if invalid.
func MustParseCIDR(raw string) CIDR {
	c, err := ParseCIDR(raw)
	if err != nil {
		panic(err)
	}
	return c
}

// IsV4 returns true if the CIDR is an IPv4 network.
func (c CIDR) IsV4() bool {
	return c.prefix.Addr().Is4()
}

// Mask returns the prefix length (bits) of the CIDR.
func (c CIDR) Mask() int {
	return c.prefix.Bits()
}

// String returns the CIDR in canonical CIDR notation.
func (c CIDR) String() string {
	return c.prefix.String()
}

// Prefix returns the underlying netip.Prefix.
func (c CIDR) Prefix() netip.Prefix {
	return c.prefix
}

// MarshalText implements encoding.TextMarshaler.
func (c CIDR) MarshalText() ([]byte, error) {
	return c.prefix.MarshalText()
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (c *CIDR) UnmarshalText(text []byte) error {
	trimmed := strings.TrimSpace(string(text))
	prefix, err := netip.ParsePrefix(trimmed)
	if err != nil {
		return err
	}
	c.prefix = prefix.Masked()
	return nil
}

// SplitTo24 splits an IPv4 CIDR into /24 SubnetBlock instances.
// It returns an explicit error if the prefix is not IPv4 or has a mask greater than 24.
func (c CIDR) SplitTo24() ([]SubnetBlock, error) {
	if !c.prefix.IsValid() {
		return nil, fmt.Errorf("invalid CIDR prefix")
	}
	if !c.IsV4() {
		return nil, fmt.Errorf("cannot split non-IPv4 CIDR %s into /24 subnets", c.prefix)
	}
	if c.prefix.Bits() > 24 {
		return nil, fmt.Errorf("CIDR %s has mask /%d greater than /24", c.prefix, c.prefix.Bits())
	}

	startIP := c.prefix.Addr()
	ipBytes := startIP.As4()
	startUint := binary.BigEndian.Uint32(ipBytes[:])
	numSubnets := 1 << (24 - c.prefix.Bits())

	subnets := make([]SubnetBlock, 0, numSubnets)
	for i := 0; i < numSubnets; i++ {
		currentUint := startUint + uint32(i)*(1<<8)
		var b [4]byte
		binary.BigEndian.PutUint32(b[:], currentUint)
		addr := netip.AddrFrom4(b)
		p := netip.PrefixFrom(addr, 24)
		subnets = append(subnets, SubnetBlock{prefix: p})
	}

	return subnets, nil
}

// SubnetBlock represents a child /24 subnet block (256 IPs).
type SubnetBlock struct {
	prefix netip.Prefix
}

// NewSubnetBlock creates a SubnetBlock from a netip.Prefix.
func NewSubnetBlock(prefix netip.Prefix) SubnetBlock {
	return SubnetBlock{prefix: prefix}
}

// Prefix returns the underlying netip.Prefix.
func (s SubnetBlock) Prefix() netip.Prefix {
	return s.prefix
}

// String returns the subnet block in canonical CIDR notation.
func (s SubnetBlock) String() string {
	return s.prefix.String()
}

// ExtractUniqueIPs extracts unique HostIPs from the subnet block.
// TODO: Implement unique IP extraction logic.
func (s SubnetBlock) ExtractUniqueIPs(count int) []HostIP {
	// TODO: Implement unique IP extraction logic.
	return nil
}

// HostIP represents a single scannable IP address.
type HostIP struct {
	addr netip.Addr
}

// NewHostIP creates a HostIP from a netip.Addr.
func NewHostIP(addr netip.Addr) HostIP {
	return HostIP{addr: addr}
}

// ParseHostIP parses an IP address string into a HostIP.
func ParseHostIP(raw string) (HostIP, error) {
	addr, err := netip.ParseAddr(strings.TrimSpace(raw))
	if err != nil {
		return HostIP{}, fmt.Errorf("invalid IP address %q: %w", raw, err)
	}
	return HostIP{addr: addr}, nil
}

// Addr returns the underlying netip.Addr.
func (h HostIP) Addr() netip.Addr {
	return h.addr
}

// String returns the string representation of the IP address.
func (h HostIP) String() string {
	return h.addr.String()
}

// WithPort returns the IP address combined with the specified port as a string (e.g. "1.1.1.1:2408" or "[2606:4700::1]:2408").
func (h HostIP) WithPort(port uint16) string {
	return netip.AddrPortFrom(h.addr, port).String()
}

// CIDRPool models a parent CIDR along with its associated child subnet blocks.
type CIDRPool struct {
	parent  CIDR
	subnets []SubnetBlock
}

// NewCIDRPool creates a new CIDRPool from a parent CIDR by internally splitting it into /24 subnets.
func NewCIDRPool(parent CIDR) (CIDRPool, error) {
	subnets, err := parent.SplitTo24()
	if err != nil {
		return CIDRPool{}, err
	}
	return CIDRPool{
		parent:  parent,
		subnets: subnets,
	}, nil
}

// Parent returns the parent CIDR of the pool.
func (p CIDRPool) Parent() CIDR {
	return p.parent
}

// Subnets returns the child SubnetBlock slice of the pool.
func (p CIDRPool) Subnets() []SubnetBlock {
	return p.subnets
}

// Extract extracts a quota of HostIPs from the pool.
// TODO: Implement quota extraction logic.
func (p *CIDRPool) Extract(quota int) []HostIP {
	// TODO: Implement quota extraction logic.
	return nil
}
