package main

import (
	"encoding/json"
	"log"
	"os"
)

type Conf struct {
	H3Addr     string `json:"H3Addr"`
	H1Addr     string `json:"H1Addr"`
	ServerName string `json:"ServerName"`
	CertPath   string `json:"CertPath"`
	KeyPath    string `json:"KeyPath"`
	Scheme     string `json:"Scheme"`
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
