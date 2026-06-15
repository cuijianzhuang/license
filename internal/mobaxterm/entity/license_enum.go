package entity

// LicenseEnum defines the type of license
type LicenseEnum int

const (
	Professional LicenseEnum = iota + 1 // Professional Edition
	Educational                         // Educational Edition
	Personal                            // Personal Edition
)

// LicenseEnumNames maps LicenseEnum values to descriptive strings
var LicenseEnumNames = map[LicenseEnum]string{
	Professional: "Professional Edition",
	Educational:  "Educational Edition",
	Personal:     "Personal Edition",
}

// GetCode returns the integer code of LicenseEnum
func (le LicenseEnum) GetCode() int {
	return int(le)
}

// GetName returns the descriptive name of LicenseEnum
func (le LicenseEnum) GetName() string {
	return LicenseEnumNames[le]
}
