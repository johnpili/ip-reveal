package models

// Config ...
type Config struct {
	HTTP struct {
		Port       int    `yaml:"port"`
		IsTLS      bool   `yaml:"is_tls"`
		ServerCert string `yaml:"server_crt"`
		ServerKey  string `yaml:"server_key"`
	} `yaml:"http"`
	Extraction struct {
		HeaderKey   string `yaml:"header_key"`
		DebugHeader bool   `yaml:"debug_header"`
	} `yaml:"extraction"`
}

// IPInfo ...
type IPInfo struct {
	IP        string `json:"ip"`
	UserAgent string `json:"user-agent"`
	IPCountry string `json:"ip-country"`
}
