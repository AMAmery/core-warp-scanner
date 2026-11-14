package generator

import (
	"fmt"
	"strings"
)

// CIDR represents a network CIDR model with a base IP prefix and subnet mask.
type CIDR struct {
	BaseIP string `json:"baseIp"`
	Mask   int    `json:"mask"`
}

// NewCIDR creates a new CIDR instance.
func NewCIDR(baseIP string, mask int) CIDR {
	return CIDR{
		BaseIP: baseIP,
		Mask:   mask,
	}
}

// String provides the full CIDR address by merging BaseIP and Mask.
// Example: BaseIP "162.159.192" and Mask 24 returns "162.159.192.0/24".
func (c CIDR) String() string {
	base := strings.TrimSpace(c.BaseIP)

	// Handle IPv4
	if strings.Contains(base, ".") {
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
	Range   CIDR          `json:"range"`
	Subnets []HostIp `json:"subnets"`
}

// NewCIDRPool creates a new CIDRPool instance.
func NewCIDRPool(parent CIDR, subnets []HostIp) CIDRPool {
	return CIDRPool{
		Range:   parent,
		Subnets: subnets,
	}
}
