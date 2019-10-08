package main

import (
  "github.com/pkg/errors"
  "encoding/base64"
  "encoding/json"
  "fmt"
  "os"
)

type jsonConfig struct {
	DefaultSecret string
  Domains []*jsonDomain
}

type jsonDomain struct {
	HostName string
	Secret string
}

type config struct {
  defaultSecret []byte
  domains map[string]*domain
}

type domain struct {
  secret []byte
}

func loadConfig(path string) (*config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
  decoder := json.NewDecoder(file)
  jsonCfg := jsonConfig{}
	err = decoder.Decode(&jsonCfg)
	if err != nil {
    return nil, err
  }

  if len(jsonCfg.DefaultSecret) < 1 {
    return nil, errors.New("DefaultSecret is not set")
  }

  defaultSecret, err := base64.StdEncoding.DecodeString(jsonCfg.DefaultSecret)
  if err != nil {
    return nil, err
  }

  domains := make(map[string]*domain)

	for _, d := range jsonCfg.Domains {
    if len(d.Secret) < 1 {
      return nil, fmt.Errorf("Secret for domain %s is not set", d.HostName)
    }

    secret, err := base64.StdEncoding.DecodeString(d.Secret)
    if err != nil {
      return nil, errors.Wrapf(err, "Domain %s", d.HostName)
    }
		domains[d.HostName] = &domain{secret}
  }

  cfg := config{defaultSecret, domains}

	return &cfg, nil
}
