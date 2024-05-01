package main

import (
	"encoding/json"
	"log"
	"os"
)

type Quic struct {
	HandshakeIdleTimeout           byte   `json:"HandshakeIdleTimeout"`
	MaxIdleTimeout                 byte   `json:"MaxIdleTimeout"`
	InitialStreamReceiveWindow     uint64 `json:"InitialStreamReceiveWindow"`
	MaxStreamReceiveWindow         uint64 `json:"MaxStreamReceiveWindow"`
	InitialConnectionReceiveWindow uint64 `json:"InitialConnectionReceiveWindow"`
	MaxConnectionReceiveWindow     uint64 `json:"MaxConnectionReceiveWindow"`
	MaxIncomingStreams             int64  `json:"MaxIncomingStreams"`
	MaxIncomingUniStreams          int64  `json:"MaxIncomingUniStreams"`
	DisablePathMTUDiscovery        bool   `json:"DisablePathMTUDiscovery"`
	Allow0RTT                      bool   `json:"Allow0RTT"`
}

type Destination struct {
	Name          string              `json:"Name"`
	Addr          string              `json:"Addr"`
	Scheme        string              `json:"Scheme"`
	Path          string              `json:"Path"`
	H3RespHeaders map[string][]string `json:"H3RespHeaders"`
	H1ReqHeaders  map[string][]string `json:"H1ReqHeaders"`
}

type Conf struct {
	H3Addr       string        `json:"H3Addr"`
	ServerName   string        `json:"ServerName"`
	CertPath     string        `json:"CertPath"`
	KeyPath      string        `json:"KeyPath"`
	Destinations []Destination `json:"Destinations"`
	QUIC         Quic          `json:"QUIC"`
}

func LoadConfig() Conf {
	cfile, cfile_err := os.ReadFile("conf.json")
	if cfile_err != nil {
		log.Fatalln(cfile_err.Error())
	}

	conf := Conf{}
	conf_err := json.Unmarshal(cfile, &conf)
	if conf_err != nil {
		log.Fatalln(conf_err.Error())
	}

	return conf
}
