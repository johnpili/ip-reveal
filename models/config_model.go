package models

// Config ...
type Config struct {
	HTTP struct {
		BasePath   string `yaml:"base_path"`
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
