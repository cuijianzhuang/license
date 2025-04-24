package sys

// Version is the version information
var Version = "0.0.1"

// Hash is the hash of commit hash
var Hash = "6eb5abee"

// Arch is the architecture of the build
var Arch = "linux/amd64"

// GetVersion returns the application version
func GetVersion() string {
	return Version
}
