package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"net"
	"strconv"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/msoedov/tcp-file/indexer"
)

var (
	portPointer        = flag.String("p", "3333", "Port")
	errMalformedPacket = errors.New("Malformed packed")
	ErrReply           = []byte("ERR\r\n")
	OkReply            = []byte("OK\r\n")
)

func main() {
	flag.Parse()
	l, err := net.Listen("tcp", ":"+*portPointer)
	if err != nil {
		log.Fatalf("Error listening: %s", err.Error())
	}
	defer l.Close()
	idx, err := indexer.BuildIndex("main.go")
	if err != nil {
		log.Fatalf("Failed to create index %v\n", err)
	}
	fmt.Println("Listening on " + "0.0.0.0:" + *portPointer)
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatalf("Error accepting: %s", err.Error())
		}
		go do(idx, conn)
	}
}

// Run run the server
func Run(hostPort string, servedFile string) {
	l, err := net.Listen("tcp", hostPort)
	if err != nil {
		log.Fatalf("Error listening: %s", err.Error())
	}
	defer l.Close()
	idx, err := indexer.BuildIndex(servedFile)
	if err != nil {
		log.Fatalf("Failed to create index %v\n", err)
	}
	log.Infof("Listening on %s\n", hostPort)
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatalf("Error accepting: %s", err.Error())
		}
		go do(idx, conn)
	}
}

func do(idx *indexer.Indexer, conn net.Conn) {
	defer conn.Close()
	ioBuffer := bufio.NewReader(conn)
	for {
		str, err := ioBuffer.ReadString('\n')
		if err != nil {
			break
		}
		if len(str) > 0 {
			line, errP := parseInt(str)
			log.WithField("payload", str).WithField("Line", line).WithField("err", errP).Info("Received")
			if errP != nil {
				conn.Write(ErrReply)
			} else {
				str, err := idx.GetLine(line)
				if err != nil {
					conn.Write(ErrReply)
					continue
				}
				conn.Write(OkReply)
				conn.Write([]byte(str))
			}
		}
	}
}

func parseInt(line string) (int64, error) {
	cmd := "GET"
	for i, char := range cmd {
		if rune(line[i]) != char {
			return -1, errMalformedPacket
		}
	}
	if line[3] != " "[0] {
		return -1, errMalformedPacket
	}
	intPart := strings.Trim(line[4:], "r\r\n\\")
	number, err := strconv.Atoi(intPart)
	if err != nil {
		return -1, err
	}
	return int64(number), nil
}
