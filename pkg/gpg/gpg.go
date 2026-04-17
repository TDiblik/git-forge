package gpg

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Manager struct {
	TempDir string
}

func NewManager() (*Manager, error) {
	tempDir, err := os.MkdirTemp("", "git-forge-gpg-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp GPG dir: %v", err)
	}

	absPath, err := filepath.Abs(tempDir)
	if err == nil {
		evalPath, err := filepath.EvalSymlinks(absPath)
		if err == nil {
			return &Manager{TempDir: evalPath}, nil
		}
		return &Manager{TempDir: absPath}, nil
	}

	return &Manager{TempDir: tempDir}, nil
}

func (m *Manager) Cleanup() {
	if m.TempDir != "" {
		os.RemoveAll(m.TempDir)
	}
}

func (m *Manager) GenerateKey(name, email string, dryRun bool) (string, error) {
	if dryRun {
		fmt.Printf("[DRY-RUN] Generating GPG key for: %s <%s> in %s\n", name, email, m.TempDir)
		return "DEADBEEF", nil
	}

	config := fmt.Sprintf(`
Key-Type: RSA
Key-Length: 2048
Subkey-Type: RSA
Subkey-Length: 2048
Name-Real: %s
Name-Email: %s
Expire-Date: 0
%%no-protection
%%commit
`, name, email)

	configFile := filepath.Join(m.TempDir, "gpg-gen-key-config")
	if err := os.WriteFile(configFile, []byte(config), 0600); err != nil {
		return "", fmt.Errorf("failed to write gpg config: %v", err)
	}

	cmd := exec.Command("gpg", "--batch", "--gen-key", configFile)
	cmd.Env = append(os.Environ(), fmt.Sprintf("GNUPGHOME=%s", m.TempDir))

	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("gpg key generation failed: %v: %s", err, string(out))
	}

	return m.extractKeyID(email)
}

func (m *Manager) extractKeyID(email string) (string, error) {
	cmd := exec.Command("gpg", "--list-keys", "--with-colons", email)
	cmd.Env = append(os.Environ(), fmt.Sprintf("GNUPGHOME=%s", m.TempDir))

	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to list keys: %v: %s", err, string(out))
	}

	lines := strings.SplitSeq(string(out), "\n")
	for line := range lines {
		if strings.HasPrefix(line, "pub:") {
			parts := strings.Split(line, ":")
			if len(parts) > 4 {
				return parts[4], nil
			}
		}
	}

	return "", fmt.Errorf("could not find key ID for email: %s", email)
}
