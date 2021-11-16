package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	port := 3333

	configFile := flag.String("config", "./srsd.json", "File containing configuration options")
	socketPath := flag.String("socket", "./srsd.sock", "Path to the socket file to listen on")
	socketAccess := flag.Uint("socket-access", 0770, "Access permissions to set for the socket file")
	flag.Parse()

	config, err := loadConfig(*configFile)
	if err != nil {
		panic(err)
  }
  
  SocketServer(port, config)

	srsd, err := newSrsd(config, *socketPath, os.FileMode(*socketAccess))
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
