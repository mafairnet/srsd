package main

import (
  "encoding/json"
  "os"
)

type config struct {
	DefaultSecret string
  Domains []*domain
  domainsByHostName map[string]*domain
}

type domain struct {
	HostName string
	Secret string
	// others - check srs.SRS
}

func loadConfig(path string) (cfg *config, err error) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()
  decoder := json.NewDecoder(file)
  cfg = &config{}
	err = decoder.Decode(cfg)
	if err != nil {
    cfg = nil
    return
  }

  domains := make(map[string]*domain)

	for _, domain := range cfg.Domains {
		domains[domain.HostName] = domain
  }

  cfg.domainsByHostName = domains

	return
}
