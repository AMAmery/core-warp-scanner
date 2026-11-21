package generator

import (
	"fmt"
	"strconv"
	"strings"
)

// CIDR represents a network CIDR model with a base IP prefix, subnet mask, and IP version indicator.
type CIDR struct {
	BaseIP string `json:"baseIp"`
	Mask   int    `json:"mask"`
	IsV4   bool   `json:"isV4"`
}

// NewCIDR creates a new CIDR instance.
func NewCIDR(baseIP string, mask int) CIDR {
	return CIDR{
		BaseIP: baseIP,
		Mask:   mask,
		IsV4:   !strings.Contains(baseIP, ":") && strings.Contains(baseIP, "."),
	}
}

// ParseCIDR creates a CIDR instance by parsing a CIDR string (e.g. "162.159.192.0/24" or "2606:4700:d1::/48").
func ParseCIDR(cidrStr string) (CIDR, error) {
	trimmed := strings.TrimSpace(cidrStr)
	parts := strings.Split(trimmed, "/")
	if len(parts) != 2 {
		return CIDR{}, fmt.Errorf("invalid CIDR format %q: expected baseIP/mask", cidrStr)
	}

	baseIP := strings.TrimSpace(parts[0])
	if baseIP == "" {
		return CIDR{}, fmt.Errorf("invalid CIDR %q: empty base IP", cidrStr)
	}

	maskStr := strings.TrimSpace(parts[1])
	mask, err := strconv.Atoi(maskStr)
	if err != nil {
		return CIDR{}, fmt.Errorf("invalid CIDR mask %q in %q: %w", maskStr, cidrStr, err)
	}

	isV4 := !strings.Contains(baseIP, ":") && strings.Contains(baseIP, ".")

	return CIDR{
		BaseIP: baseIP,
		Mask:   mask,
		IsV4:   isV4,
	}, nil
}

// NewCIDRFromString creates a new CIDR instance from a CIDR string (alias to ParseCIDR).
func NewCIDRFromString(cidrStr string) (CIDR, error) {
	return ParseCIDR(cidrStr)
}

// MustParseCIDR creates a new CIDR instance from a CIDR string and panics if invalid.
func MustParseCIDR(cidrStr string) CIDR {
	c, err := ParseCIDR(cidrStr)
	if err != nil {
		panic(err)
	}
	return c
}

// String provides the full CIDR address by merging BaseIP and Mask.
// Example: BaseIP "162.159.192" and Mask 24 returns "162.159.192.0/24".
func (c CIDR) String() string {
	base := strings.TrimSpace(c.BaseIP)

	// Handle IPv4
	if c.IsV4 || strings.Contains(base, ".") {
		parts := strings.Split(base, ".")
		for len(parts) < 4 {
			base += ".0"
			parts = append(parts, "0")
		}
		return fmt.Sprintf("%s/%d", base, c.Mask)
	}

	// Handle IPv6
	if strings.Contains(base, ":") {
		if !strings.HasSuffix(base, "::") {
			if strings.HasSuffix(base, ":") {
				base += ":"
			} else {
				base += "::"
			}
		}
		return fmt.Sprintf("%s/%d", base, c.Mask)
	}

	return fmt.Sprintf("%s/%d", base, c.Mask)
}

// HostIp represents an individual IP address block within a subnet.
type HostIp struct {
	IP string `json:"ip"`
}

// NewHostIp creates a new HostIp instance.
func NewHostIp(ip string) HostIp {
	return HostIp{
		IP: ip,
	}
}

// CIDRPool models a parent CIDR along with its associated child subnet blocks.
type CIDRPool struct {
	Range   CIDR     `json:"range"`
	Subnets []HostIp `json:"subnets"`
}

// NewCIDRPool creates a new CIDRPool instance.
func NewCIDRPool(parent CIDR, subnets []HostIp) CIDRPool {
	return CIDRPool{
		Range:   parent,
		Subnets: subnets,
	}
}
