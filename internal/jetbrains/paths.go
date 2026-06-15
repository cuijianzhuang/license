package jetbrains

import (
	"license/internal/config"
	"path/filepath"
)

// Paths holds the on-disk locations of the JetBrains certificate material.
// Keeping these on FakeCert (instead of package-level vars) means tests and
// alternative deployments can supply their own paths without mutating shared
// state.
type Paths struct {
	CodeRootCert   string
	ServerRootCert string
	PrivateKey     string
	PublicKey      string
	CodeCert       string
	ServerCert     string
}

// DefaultPaths returns paths rooted at the configured data directory. Used
// when no certificate manager is available (legacy/backward-compatible mode).
func DefaultPaths() Paths {
	dataDir := config.GetConfig().DataDir
	return Paths{
		CodeRootCert:   filepath.Join(dataDir, "jetbrainsCodeCACert.pem"),
		ServerRootCert: filepath.Join(dataDir, "jetbrainsServerCACert.pem"),
		PrivateKey:     filepath.Join(dataDir, "private.pem"),
		PublicKey:      filepath.Join(dataDir, "public.pem"),
		CodeCert:       filepath.Join(dataDir, "code.pem"),
		ServerCert:     filepath.Join(dataDir, "server.pem"),
	}
}

// certManager is the minimum surface the initialize package's CertManager
// must satisfy. Declaring it here avoids an import cycle.
type certManager interface {
	GetFilePath(string) string
}

// PathsFromCertManager builds Paths by asking a certificate manager for the
// canonical location of each named artifact.
func PathsFromCertManager(cm certManager) Paths {
	return Paths{
		CodeRootCert:   cm.GetFilePath("jetbrains_code_ca"),
		ServerRootCert: cm.GetFilePath("jetbrains_server_ca"),
		PrivateKey:     cm.GetFilePath("jetbrains_private_key"),
		PublicKey:      cm.GetFilePath("jetbrains_public_key"),
		CodeCert:       cm.GetFilePath("jetbrains_code_cert"),
		ServerCert:     cm.GetFilePath("jetbrains_server_cert"),
	}
}
