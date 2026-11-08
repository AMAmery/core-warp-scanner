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

func TestAssetIPProvider_GetIPv4_Default(t *testing.T) {
	provider := &AssetIPProvider{}

	ipv4s, err := provider.GetIPv4()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(ipv4s) == 0 {
		t.Fatal("expected non-empty list of IPv4 addresses")
	}

	// Verify some known contents from assets/ipv4.txt
	expectedFirst := "173.245.48.0/20"
	if ipv4s[0] != expectedFirst {
		t.Errorf("expected first IPv4 to be %q, got %q", expectedFirst, ipv4s[0])
	}

	if len(ipv4s) != 15 {
		t.Errorf("expected 15 IPv4 ranges, got %d", len(ipv4s))
	}
}

func TestAssetIPProvider_GetIPv6_Default(t *testing.T) {
	provider := &AssetIPProvider{}

	ipv6s, err := provider.GetIPv6()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(ipv6s) == 0 {
		t.Fatal("expected non-empty list of IPv6 addresses")
	}

	// Verify some known contents from assets/ipv6.txt
	expectedFirst := "2400:cb00::/32"
	if ipv6s[0] != expectedFirst {
		t.Errorf("expected first IPv6 to be %q, got %q", expectedFirst, ipv6s[0])
	}

	if len(ipv6s) != 7 {
		t.Errorf("expected 7 IPv6 ranges, got %d", len(ipv6s))
	}
}

func TestAssetIPProvider_CustomAssetPaths(t *testing.T) {
	tempDir := t.TempDir()

	customV4Path := filepath.Join(tempDir, "custom_v4.txt")
	customV6Path := filepath.Join(tempDir, "custom_v6.txt")

	v4Content := "# Comment\n1.1.1.0/24\n\n1.0.0.0/24\n"
	v6Content := "2606:4700::/32\n"

	if err := os.WriteFile(customV4Path, []byte(v4Content), 0644); err != nil {
		t.Fatalf("failed to write custom v4 test file: %v", err)
	}
	if err := os.WriteFile(customV6Path, []byte(v6Content), 0644); err != nil {
		t.Fatalf("failed to write custom v6 test file: %v", err)
	}

	// Temporarily override assetPaths struct
	orig := assetPaths
	defer func() { assetPaths = orig }()

	assetPaths.ipv4 = customV4Path
	assetPaths.ipv6 = customV6Path

	provider := &AssetIPProvider{}

	v4s, err := provider.GetIPv4()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(v4s) != 2 || v4s[0] != "1.1.1.0/24" || v4s[1] != "1.0.0.0/24" {
		t.Errorf("unexpected v4s: %#v", v4s)
	}

	v6s, err := provider.GetIPv6()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(v6s) != 1 || v6s[0] != "2606:4700::/32" {
		t.Errorf("unexpected v6s: %#v", v6s)
	}
}

func TestAssetIPProvider_FileNotFound(t *testing.T) {
	orig := assetPaths
	defer func() { assetPaths = orig }()

	assetPaths.ipv4 = "non_existent_v4.txt"
	assetPaths.ipv6 = "non_existent_v6.txt"

	provider := &AssetIPProvider{}

	_, errV4 := provider.GetIPv4()
	if errV4 == nil {
		t.Error("expected error for non-existent IPv4 file, got nil")
	}

	_, errV6 := provider.GetIPv6()
	if errV6 == nil {
		t.Error("expected error for non-existent IPv6 file, got nil")
	}
}
