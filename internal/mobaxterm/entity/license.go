package entity

import (
	"fmt"
)

// License defines the detailed information for a license
type License struct {
	LicenseType  LicenseEnum
	UserName     string
	MajorVersion int
	MinorVersion int
	Count        int
	Unknown      int
	OpenGames    bool
	OpenPlugins  bool
}

// NewLicense creates a new License instance
func NewLicense(licenseType LicenseEnum, userName string, majorVersion, minorVersion, count int, openGames, openPlugins bool) *License {
	return &License{
		LicenseType:  licenseType,
		UserName:     userName,
		MajorVersion: majorVersion,
		MinorVersion: minorVersion,
		Count:        count,
		Unknown:      0, // Default value
		OpenGames:    openGames,
		OpenPlugins:  openPlugins,
	}
}

// GetLicenseKey generates the string for the key to be calculated
func (l *License) GetLicenseKey() string {
	game := 0
	if l.OpenGames {
		game = 1
	}
	plugin := 0
	if l.OpenPlugins {
		plugin = 1
	}
	// Use the same string format that matches the format in Java
	licenseFormat := "%d#%s|%d%d#%d#%d3%d6%d#%d#%d#%d#"
	return fmt.Sprintf(licenseFormat, int(l.LicenseType), l.UserName, l.MajorVersion, l.MinorVersion, l.Count, l.MajorVersion, l.MinorVersion, l.MinorVersion, l.Unknown, game, plugin)
}
