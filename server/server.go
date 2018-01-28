package server

import (
	"bufio"
	"errors"
	"net"
	"strconv"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/msoedov/tcp-file/indexer"
)

var (
	errMalformedPacket = errors.New("Malformed packed")
	errQuit            = errors.New("Quit cmd")
	ErrReply           = []byte("ERR\r\n")
	OkReply            = []byte("OK\r\n")
	CmdQuit            = "QUIT"
	CmdGet             = "GET"
)

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
			if errP == errQuit {
				conn.Write(OkReply)
				return
			}
			log.WithField("payload", str).WithField("Line", line).WithField("err", errP).Debug("Received")
			if errP != nil {
				conn.Write(ErrReply)
			} else {
				payload, err := idx.GetLineBytes(line)
				if err != nil {
					conn.Write(ErrReply)
					continue
				}
				// Avoid extra mem alocation
				conn.Write(append(OkReply, payload...))
			}
		}
	}
}

func parseInt(line string) (int64, error) {
	if len(line) < 4 {
		return -1, errMalformedPacket
	}
	if line[:4] == CmdQuit {
		return -1, errQuit
	}
	if line[:3] != CmdGet {
		return -1, errMalformedPacket
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
