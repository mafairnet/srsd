package main

import (
	"bufio"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

var configuration = getProgramConfiguration()

func SocketServer(port int, config *config) {

	listen, err := net.Listen("tcp4", ":"+strconv.Itoa(port))

	if err != nil {
		log.Fatalf("Socket listen port %d failed,%s", port, err)
		os.Exit(1)
	}

	defer listen.Close()

	log.Printf("Begin listen port: %d", port)

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Fatalln(err)
			continue
		}
		go handler(conn, config)
	}

}
func handler(conn net.Conn, config *config) {

	defer conn.Close()

	var (
		buf = make([]byte, 1024)
		r   = bufio.NewReader(conn)
		w   = bufio.NewWriter(conn)
	)

ILOOP:
	for {
		n, err := r.Read(buf)
		data := string(buf[:n])

		switch err {
		case io.EOF:
			break ILOOP
		case nil:

			log.Println("Receive:", data)
			data = strings.TrimSuffix(data, "\n")
			cmd, err := parseCommand(data)
			log.Printf("Command Result: %s", cmd)

			if err != nil {
				//writeError(writer, err)
				errorString := err.Error()
				log.Printf("Error: %s", errorString)
				w.Write([]byte(errorString))
				w.Write([]byte(configuration.StopCharacter))
				break ILOOP
			}

			result, err := cmd.process(config)

			if err != nil {
				errorString := err.Error()
				log.Printf("Error: %s", errorString)
				w.Write([]byte(errorString))
				w.Write([]byte(configuration.StopCharacter))
			} else {
				result = "200 " + result
				w.Write([]byte(result))
				w.Write([]byte(configuration.StopCharacter))
				w.Flush()
				log.Printf("Sent: %s", result)
			}

			if isTransportOver(data) {
				break ILOOP
			}

		default:
			log.Fatalf("Receive data failed:%s", err)
			return
		}

	}

}

func isTransportOver(data string) (over bool) {
	over = strings.HasSuffix(data, "\n")
	return
}
