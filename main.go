package main

import (
  "flag"
  "os"
  "os/signal"
  "syscall"
)

func main() {
  configFile := flag.String("config", "./srsd.json", "File containing configuration options")
  socketPath := flag.String("socket", "./srsd.sock", "Path to the socket file to listen on")
	flag.Parse()

  config, err := loadConfig(*configFile)
	if err != nil {
		panic(err)
	}

  srsd, err := newSrsd(config, *socketPath)
	if err != nil {
		panic(err)
	}

  signals := make(chan os.Signal, 1)
  signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
  go func() {
    for _ = range signals {
      srsd.close()
      os.Exit(1)
    }
  }()

  defer srsd.close()
  srsd.run()
}
