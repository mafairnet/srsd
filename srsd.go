package main

import (
  "bufio"
  "net"
  "os"
  "strings"
)

type srsdListener struct {
  config *config
  socketPath string
  listener net.Listener
}

type srsd interface {
  run() error
  close() error
}

func newSrsd(config *config, socketPath string) (srsd srsd, err error) {
  os.Remove(socketPath)

  listener, err := net.Listen("unix", socketPath)
  if err != nil {
    return
  }

	srsd = &srsdListener{config, socketPath, listener}
  return
}

func (srsd *srsdListener) run() (error) {
  for {
		connection, err := srsd.listener.Accept()
		if err != nil {
			return err
		}

		go handleConnection(srsd, connection)
  }

  return nil
}

func writeError(writer *bufio.Writer, writeErr error) (err error) {
  _, err = writer.WriteString("ERROR: ")
  if err != nil {
    return
  }
  _, err = writer.WriteString(writeErr.Error())
  if err != nil {
    return
  }

  _, err = writer.WriteString("\n")
  if err != nil {
    return
  }

  err = writer.Flush()
  return
}

func handleConnection(srsd *srsdListener, connection net.Conn) {
	reader := bufio.NewReader(connection)
  writer := bufio.NewWriter(connection)

  for {
    line, err := reader.ReadString('\n')
    if err != nil {
      break;
    }

    line = strings.TrimSuffix(line, "\n")

    cmd, err := parseCommand(line)
    if err != nil {
      writeError(writer, err)
      break;
    }

    result, err := cmd.process(srsd.config)
    if err != nil {
      writeError(writer, err)
    } else {
      _, err = writer.WriteString(result)
      if err != nil {
        break;
      }

      _, err = writer.WriteString("\n")
      if err != nil {
        break;
      }
    }

    err = writer.Flush()
    if err != nil {
      break;
    }
  }

	connection.Close()
}

func (srsd *srsdListener) close() (err error) {
  err = srsd.listener.Close()
  if err != nil {
    return
  }

  err = os.Remove(srsd.socketPath)
	return
}
