package main

import (
  "github.com/mileusna/srs"
  "errors"
  "strings"
)

type forwardCommand struct {
  srsDomain string
  emailAddress string
}

type reverseCommand struct {
  emailAddress string
}

type command interface {
  process(*config) (string, error)
}

func parseCommand(s string) (command, error) {
  if len(s) < 1 {
    return nil, errors.New("Empty command")
  }

  if s[0] == 'F' {
    parts := strings.SplitN(s, ":", 2)

    if len(parts) != 2 {
      return nil, errors.New("Missing SRS domain or email address")
    }

    srsDomain := parts[0][1:]
    return &forwardCommand{srsDomain, parts[1]}, nil
  } else if s[0] == 'R' {
    emailAddress := s[1:]
    return &reverseCommand{emailAddress}, nil
  }

  return nil, errors.New("Invalid command type")
}

func newSrs(cfg *config, srsDomain string) srs.SRS {
  srs := srs.SRS{Domain: srsDomain}

  domain := cfg.domains[srsDomain]
  if domain == nil {
    srs.Secret = cfg.defaultSecret
  } else {
    srs.Secret = domain.secret
  }

  return srs
}

func (cmd *forwardCommand) process(cfg *config) (string, error) {
  srs := newSrs(cfg, cmd.srsDomain)
  return srs.Forward(cmd.emailAddress)
}

func (cmd *reverseCommand) process(cfg *config) (string, error) {
  parts := strings.Split(cmd.emailAddress, "@")
  if len(parts) < 1 {
    return "", errors.New("Missing domain in email address")
  }
  srsDomain := parts[len(parts) - 1]
  srs := newSrs(cfg, srsDomain)
  return srs.Reverse(cmd.emailAddress)
}
