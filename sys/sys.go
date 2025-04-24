package sys

// Version is the version information
var Version = "0.0.1"

// Hash is the hash of build
var Hash = "2ed26fe1"

// Arch is the architecture of the build
var Arch = "linux/amd64"

// GetVersion returns the application version
func GetVersion() string {
	return Version
}

// GetHash returns the build hash
func GetHash() string {
	return Hash
}

// GetArch returns the build architecture
func GetArch() string {
	return Arch
}
