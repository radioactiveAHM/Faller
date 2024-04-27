package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Conf struct {
	H3Addr     string `json:"H3Addr"`
	H1Addr     string `json:"H1Addr"`
	ServerName string `json:"ServerName"`
	CertPath   string `json:"CertPath"`
	KeyPath    string `json:"KeyPath"`
}

func LoadConfig() Conf {
	cfile, cfile_err := os.ReadFile("conf.json")
	if cfile_err != nil {
		fmt.Println(cfile_err.Error())
		os.Exit(1)
	}

	conf := Conf{}
	conf_err := json.Unmarshal(cfile, &conf)
	if conf_err != nil {
		fmt.Println(conf_err.Error())
		os.Exit(1)
	}

	return conf
}
