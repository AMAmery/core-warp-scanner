package generator

import (
	"os"
	"path/filepath"
	"testing"
)

func TestAssetIPProvider_InterfaceImplementation(t *testing.T) {
	var _ IPProvider = &AssetIPProvider{}
	var _ IPProvider = AssetIPProvider{}
}

func TestAssetIPProvider_GetIPv4(t *testing.T) {
	provider := AssetIPProvider{}

	cidrs, err := provider.GetIPv4()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(cidrs) == 0 {
		t.Fatal("expected non-empty list of IPv4 CIDR ranges")
	}

	if len(cidrs) != 7 {
		t.Errorf("expected 7 IPv4 ranges, got %d", len(cidrs))
	}

	// Verify first entry from assets/ipv4.txt (162.159.192.0/24)
	first := cidrs[0]
	if first.BaseIP != "162.159.192.0" {
		t.Errorf("expected first BaseIP to be %q, got %q", "162.159.192.0", first.BaseIP)
	}
	if first.Mask != 24 {
		t.Errorf("expected first Mask to be 24, got %d", first.Mask)
	}
	if !first.IsV4 {
		t.Errorf("expected first IsV4 to be true, got %v", first.IsV4)
	}

	for i, c := range cidrs {
		if !c.IsV4 {
			t.Errorf("entry %d (%s) expected IsV4 to be true", i, c.BaseIP)
		}
		if c.Mask <= 0 {
			t.Errorf("entry %d (%s) expected positive Mask, got %d", i, c.BaseIP, c.Mask)
		}
	}
}

func TestAssetIPProvider_GetIPv6(t *testing.T) {
	provider := AssetIPProvider{}

	cidrs, err := provider.GetIPv6()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(cidrs) == 0 {
		t.Fatal("expected non-empty list of IPv6 CIDR ranges")
	}

	if len(cidrs) != 2 {
		t.Errorf("expected 2 IPv6 ranges, got %d", len(cidrs))
	}

	// Verify first entry from assets/ipv6.txt (2606:4700:d0::/48)
	first := cidrs[0]
	if first.BaseIP != "2606:4700:d0::" {
		t.Errorf("expected first BaseIP to be %q, got %q", "2606:4700:d0::", first.BaseIP)
	}
	if first.Mask != 48 {
		t.Errorf("expected first Mask to be 48, got %d", first.Mask)
	}
	if first.IsV4 {
		t.Errorf("expected first IsV4 to be false, got %v", first.IsV4)
	}

	for i, c := range cidrs {
		if c.IsV4 {
			t.Errorf("entry %d (%s) expected IsV4 to be false", i, c.BaseIP)
		}
		if c.Mask <= 0 {
			t.Errorf("entry %d (%s) expected positive Mask, got %d", i, c.BaseIP, c.Mask)
		}
	}
}

func TestAssetIPProvider_CustomAssetPaths(t *testing.T) {
	tempDir := t.TempDir()

	customV4Path := filepath.Join(tempDir, "custom_v4.txt")
	customV6Path := filepath.Join(tempDir, "custom_v6.txt")

	v4Content := "# Comment\n1.1.1.0/24\n\n1.0.0.0/24\n"
	v6Content := "# IPv6 Comment\n2606:4700::/32\n"

	if err := os.WriteFile(customV4Path, []byte(v4Content), 0644); err != nil {
		t.Fatalf("failed to write custom v4 test file: %v", err)
	}
	if err := os.WriteFile(customV6Path, []byte(v6Content), 0644); err != nil {
		t.Fatalf("failed to write custom v6 test file: %v", err)
	}

	// Temporarily override assetPaths
	orig := assetPaths
	defer func() { assetPaths = orig }()

	assetPaths.ipv4 = customV4Path
	assetPaths.ipv6 = customV6Path

	provider := AssetIPProvider{}

	v4s, err := provider.GetIPv4()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(v4s) != 2 {
		t.Fatalf("expected 2 v4 entries, got %d", len(v4s))
	}
	if v4s[0].BaseIP != "1.1.1.0" || v4s[0].Mask != 24 || !v4s[0].IsV4 {
		t.Errorf("unexpected first v4 entry: %+v", v4s[0])
	}
	if v4s[1].BaseIP != "1.0.0.0" || v4s[1].Mask != 24 || !v4s[1].IsV4 {
		t.Errorf("unexpected second v4 entry: %+v", v4s[1])
	}

	v6s, err := provider.GetIPv6()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(v6s) != 1 {
		t.Fatalf("expected 1 v6 entry, got %d", len(v6s))
	}
	if v6s[0].BaseIP != "2606:4700::" || v6s[0].Mask != 32 || v6s[0].IsV4 {
		t.Errorf("unexpected v6 entry: %+v", v6s[0])
	}
}

func TestAssetIPProvider_FileNotFound(t *testing.T) {
	orig := assetPaths
	defer func() { assetPaths = orig }()

	assetPaths.ipv4 = "non_existent_v4.txt"
	assetPaths.ipv6 = "non_existent_v6.txt"

	provider := AssetIPProvider{}

	if _, err := provider.GetIPv4(); err == nil {
		t.Error("expected error for non-existent IPv4 file, got nil")
	}

	if _, err := provider.GetIPv6(); err == nil {
		t.Error("expected error for non-existent IPv6 file, got nil")
	}
}

func TestAssetIPProvider_InvalidCIDRFile(t *testing.T) {
	tempDir := t.TempDir()
	invalidFilePath := filepath.Join(tempDir, "invalid.txt")

	if err := os.WriteFile(invalidFilePath, []byte("invalid-cidr-line\n"), 0644); err != nil {
		t.Fatalf("failed to write invalid test file: %v", err)
	}

	orig := assetPaths
	defer func() { assetPaths = orig }()

	assetPaths.ipv4 = invalidFilePath

	provider := AssetIPProvider{}
	if _, err := provider.GetIPv4(); err == nil {
		t.Error("expected error parsing invalid CIDR from file, got nil")
	}
}

