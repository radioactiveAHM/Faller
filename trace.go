package main

import (
	"log"

	"github.com/fatih/color"
)

func TraceLogReq(method string, path string, proto string, remoteAddr string, id int) {
	log.Println(color.CyanString("%s %s %s %s (ID=%d)", method, path, proto, remoteAddr, id))
}

func TraceLogResp(proto string, statuscode int, id int) {
	log.Println(color.GreenString("%s %d (ID=%d)", proto, statuscode, id))
}
