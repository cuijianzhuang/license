package initialize

import "license/config"

// ExecuteInitialize initializes various components
func ExecuteInitialize() {
	dataDir := config.GetConfig().DataDir

	// Initialize certificates
	InitCert(dataDir)
	// Initialize GitLab
	InitGitLabCert()
	// Initialize JetBrains
	InitJetbrains()
}
