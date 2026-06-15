package initialize

import (
	"fmt"
	"os"
	"path/filepath"
)

// CertFile represents a certificate file with its configuration
type CertFile struct {
	Name    string
	Path    string
	Content string
	Perm    os.FileMode
}

// CertManager manages all certificate files. The files map is built once by
// NewCertManager and never mutated afterwards, so no locking is needed.
type CertManager struct {
	dataDir string
	files   map[string]*CertFile
}

// Certificate path constants
const (
	GitLabPrivateKey    = "gitlab_private_key"
	GitLabPublicKey     = "gitlab_public_key"
	JetBrainsCodeCA     = "jetbrains_code_ca"
	JetBrainsServerCA   = "jetbrains_server_ca"
	JetBrainsPrivateKey = "jetbrains_private_key"
	JetBrainsPublicKey  = "jetbrains_public_key"
	JetBrainsCodeCert   = "jetbrains_code_cert"
	JetBrainsServerCert = "jetbrains_server_cert"
)

// NewCertManager creates a new certificate manager
func NewCertManager(dataDir string) *CertManager {
	cm := &CertManager{
		dataDir: dataDir,
		files:   make(map[string]*CertFile),
	}
	cm.initializeCertFiles()
	return cm
}

// initializeCertFiles sets up all certificate file configurations
func (cm *CertManager) initializeCertFiles() {
	cm.files = map[string]*CertFile{
		GitLabPrivateKey: {
			Name:    "GitLab Private Key",
			Path:    ".license_decryption_key.pri",
			Content: gitlabPrivateKey,
			Perm:    0600,
		},
		GitLabPublicKey: {
			Name:    "GitLab Public Key",
			Path:    ".license_encryption_key.pub",
			Content: gitlabPublicKey,
			Perm:    0600,
		},
		JetBrainsCodeCA: {
			Name:    "JetBrains Code CA",
			Path:    "jetbrainsCodeCACert.pem",
			Content: jetbrainsCodeCa,
			Perm:    0644,
		},
		JetBrainsServerCA: {
			Name:    "JetBrains Server CA",
			Path:    "jetbrainsServerCACert.pem",
			Content: jetbrainsServerCa,
			Perm:    0644,
		},
		JetBrainsPrivateKey: {
			Name: "JetBrains Private Key",
			Path: "private.pem",
			Perm: 0600,
		},
		JetBrainsPublicKey: {
			Name: "JetBrains Public Key",
			Path: "public.pem",
			Perm: 0600,
		},
		JetBrainsCodeCert: {
			Name: "JetBrains Code Certificate",
			Path: "code.pem",
			Perm: 0644,
		},
		JetBrainsServerCert: {
			Name: "JetBrains Server Certificate",
			Path: "server.pem",
			Perm: 0644,
		},
	}
}

// GetFilePath returns the full path for a certificate file
func (cm *CertManager) GetFilePath(certType string) string {
	if cert, ok := cm.files[certType]; ok {
		return filepath.Join(cm.dataDir, cert.Path)
	}
	return ""
}

// EnsureFile ensures a certificate file exists with the correct content
func (cm *CertManager) EnsureFile(certType string) error {
	cert, ok := cm.files[certType]
	if !ok {
		return fmt.Errorf("unknown certificate type: %s", certType)
	}

	// Skip if no default content is provided
	if cert.Content == "" {
		return nil
	}

	fullPath := filepath.Join(cm.dataDir, cert.Path)

	// Check if file exists
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		// Create directory if it doesn't exist
		dir := filepath.Dir(fullPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}

		// Write file with specified permissions
		if err := os.WriteFile(fullPath, []byte(cert.Content), cert.Perm); err != nil {
			return fmt.Errorf("failed to write %s: %w", cert.Name, err)
		}
	}

	return nil
}

// InitializeAll initializes all certificate files
func (cm *CertManager) InitializeAll() error {
	// Initialize basic certificates (with content)
	basicCerts := []string{
		GitLabPrivateKey,
		GitLabPublicKey,
		JetBrainsCodeCA,
		JetBrainsServerCA,
	}

	for _, certType := range basicCerts {
		if err := cm.EnsureFile(certType); err != nil {
			return fmt.Errorf("failed to initialize %s: %w", certType, err)
		}
	}

	return nil
}

// FileExists checks if a certificate file exists
func (cm *CertManager) FileExists(certType string) bool {
	path := cm.GetFilePath(certType)
	if path == "" {
		return false
	}

	_, err := os.Stat(path)
	return err == nil
}

// ReadFile reads the content of a certificate file
func (cm *CertManager) ReadFile(certType string) ([]byte, error) {
	path := cm.GetFilePath(certType)
	if path == "" {
		return nil, fmt.Errorf("unknown certificate type: %s", certType)
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s: %w", certType, err)
	}

	return content, nil
}

// WriteFile writes content to a certificate file
func (cm *CertManager) WriteFile(certType string, content []byte) error {
	cert, ok := cm.files[certType]
	if !ok {
		return fmt.Errorf("unknown certificate type: %s", certType)
	}

	fullPath := filepath.Join(cm.dataDir, cert.Path)

	// Create directory if it doesn't exist
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	// Write file with specified permissions
	if err := os.WriteFile(fullPath, content, cert.Perm); err != nil {
		return fmt.Errorf("failed to write %s: %w", cert.Name, err)
	}

	return nil
}

// GetAllPaths returns all certificate paths for initialization
func (cm *CertManager) GetAllPaths() map[string]string {
	paths := make(map[string]string)
	for key, cert := range cm.files {
		paths[key] = filepath.Join(cm.dataDir, cert.Path)
	}
	return paths
}
