package models

// Config ...
type Config struct {
	HTTP struct {
		Port       int    `yaml:"port"`
		ServerCert string `yaml:"server_crt"`
		ServerKey  string `yaml:"server_key"`
	} `yaml:"http"`
	Extraction struct {
		HeaderKey string `yaml:"header_key"`
	} `yaml:"extraction"`
}

// IPInfo ...
type IPInfo struct {
	IP        string `json:"ip"`
	UserAgent string `json:"user-agent"`
}
