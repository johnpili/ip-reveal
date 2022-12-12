package models

// IPInfo ...
type IPInfo struct {
	IP        string `json:"ip"`
	UserAgent string `json:"user-agent"`
	IPCountry string `json:"ip-country"`
}
