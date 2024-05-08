package main

import (
	"errors"
	"strings"
)

type Domainy struct {
	Subdomain  string
	SDL        string
	TLD        string
	DomainName string
}

func (d *Domainy) From(domain string) error {

	stat := true

	parts := strings.Split(domain, ".")
	if safetyCheck(parts) {
		stat = false
		return errors.New("invalid domain")
	}

	defer func() {
		if stat {
			d.DomainName = strings.ReplaceAll(d.DomainName, " ", "")
			d.SDL = strings.ReplaceAll(d.SDL, " ", "")
			d.Subdomain = strings.ReplaceAll(d.Subdomain, " ", "")
			d.TLD = strings.ReplaceAll(d.TLD, " ", "")
		}
	}()

	dlen := len(parts)
	if dlen == 1 {
		d.SDL = parts[0]
		return nil
	} else if dlen == 2 {
		d.SDL = parts[0]
		d.TLD = "." + parts[1]
		d.DomainName = d.SDL + d.TLD
		return nil
	} else if dlen > 2 {
		d.SDL = parts[dlen-2]
		d.TLD = "." + parts[dlen-1]
		d.DomainName = d.SDL + d.TLD
		d.Subdomain = strings.Join(parts[:dlen-2], ".")
		return nil
	}

	return errors.New("failed to parse domain")
}

func safetyCheck(l []string) bool {
	for _, element := range l {
		if element == "" || element == " " || element == "." || len(element) == 0 {
			return true
		}
	}

	return false
}
