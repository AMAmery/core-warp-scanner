package generator

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// IPProvider defines the abstraction for providing IPv4 and IPv6 addresses.
type IPProvider interface {
	GetIPv4() ([]string, error)
	GetIPv6() ([]string, error)
}

// assetPaths defines the asset file paths as a private struct (similar to an object in Kotlin).
var assetPaths = struct {
	ipv4 string
	ipv6 string
}{
	ipv4: "assets/ipv4.txt",
	ipv6: "assets/ipv6.txt",
}

// AssetIPProvider implements IPProvider by reading IP addresses from text asset files.
type AssetIPProvider struct{}

// Ensure AssetIPProvider implements IPProvider at compile time.
var _ IPProvider = (*AssetIPProvider)(nil)
var _ IPProvider = AssetIPProvider{}

// GetIPv4 reads and returns IPv4 addresses/prefixes from the configured IPv4 asset file.
func (AssetIPProvider) GetIPv4() ([]string, error) {
	return readLines(assetPaths.ipv4)
}

// GetIPv6 reads and returns IPv6 addresses/prefixes from the configured IPv6 asset file.
func (AssetIPProvider) GetIPv6() ([]string, error) {
	return readLines(assetPaths.ipv6)
}

// readLines reads a file line by line, trimming whitespace and ignoring empty lines.
func readLines(path string) ([]string, error) {
	resolvedPath := resolvePath(path)
	file, err := os.Open(resolvedPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open asset file %q: %w", path, err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		lines = append(lines, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading asset file %q: %w", path, err)
	}

	return lines, nil
}

// resolvePath checks if path exists directly; if not, checks parent directory (useful when running tests from subpackages).
func resolvePath(path string) string {
	if _, err := os.Stat(path); err == nil {
		return path
	}
	parentPath := filepath.Join("..", path)
	if _, err := os.Stat(parentPath); err == nil {
		return parentPath
	}
	return path
}
