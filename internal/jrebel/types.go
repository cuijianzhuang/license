package jrebel

// LeasesResponse is the JSON body returned for POST /leases.
type LeasesResponse struct {
	ServerVersion         string   `json:"serverVersion"`
	ServerProtocolVersion string   `json:"serverProtocolVersion"`
	ServerGuid            string   `json:"serverGuid"`
	GroupType             string   `json:"groupType"`
	ID                    int64    `json:"id"`
	LicenseType           int64    `json:"licenseType"`
	EvaluationLicense     bool     `json:"evaluationLicense"`
	Signature             string   `json:"signature"`
	ServerRandomness      string   `json:"serverRandomness"`
	SeatPoolType          string   `json:"seatPoolType"`
	StatusCode            string   `json:"statusCode"`
	Offline               bool     `json:"offline"`
	ValidFrom             int64    `json:"validFrom"`
	ValidUntil            int64    `json:"validUntil"`
	Company               string   `json:"company"`
	OrderId               string   `json:"orderId"`
	ZeroIds               []string `json:"zeroIds"`
	LicenseValidFrom      int64    `json:"licenseValidFrom"`
	LicenseValidUntil     int64    `json:"licenseValidUntil"`
}

// LeasesOneResponse is the JSON body returned for DELETE /leases/1.
type LeasesOneResponse struct {
	ServerVersion         string `json:"serverVersion"`
	ServerProtocolVersion string `json:"serverProtocolVersion"`
	ServerGuid            string `json:"serverGuid"`
	GroupType             string `json:"groupType"`
	StatusCode            string `json:"statusCode"`
	Signature             string `json:"signature"`
	ServerRandomness      string `json:"serverRandomness"`
	Features              string `json:"features"`
	Msg                   string `json:"msg"`
	StatusMessage         string `json:"statusMessage"`
	Company               string `json:"company"`
}

// ValidateResponse is the JSON body returned for POST /validate-connection and
// /features.
type ValidateResponse struct {
	ServerVersion         string `json:"serverVersion"`
	ServerProtocolVersion string `json:"serverProtocolVersion"`
	ServerGuid            string `json:"serverGuid"`
	GroupType             string `json:"groupType"`
	StatusCode            string `json:"statusCode"`
	Signature             string `json:"signature"`
	ServerRandomness      string `json:"serverRandomness"`
	Features              string `json:"features"`
	Company               string `json:"company"`
	CanGetLease           bool   `json:"canGetLease"`
	LicenseType           string `json:"licenseType"`
	EvaluationLicense     bool   `json:"evaluationLicense"`
	SeatPoolType          string `json:"seatPoolType"`
}
