package main

type Domains struct {
	ServerName        string `json:"ServerName"`
	CertPath          string `json:"CertPath"`
	KeyPath           string `json:"KeyPath"`
	SubDomainsSupport bool   `json:"SubDomainsSupport"`
}

type DefaultTls struct {
	CertPath string `json:"CertPath"`
	KeyPath  string `json:"KeyPath"`
}

type TlsConf struct {
	Default DefaultTls `json:"Default"`
	Domains []Domains  `json:"Domains"`
}
