package entity

// LicenseEnum is the MobaXterm edition code embedded in the license payload.
type LicenseEnum int

const (
	Professional LicenseEnum = iota + 1 // Professional Edition
	Educational                         // Educational Edition
	Personal                            // Personal Edition
)

// String returns the human-readable edition name.
func (le LicenseEnum) String() string {
	switch le {
	case Professional:
		return "Professional Edition"
	case Educational:
		return "Educational Edition"
	case Personal:
		return "Personal Edition"
	}
	return ""
}
