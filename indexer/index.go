package indexer

import (
	"bufio"
	"encoding/binary"
	"io"
	"os"
	"sync"

	log "github.com/Sirupsen/logrus"
)

// Indexer documents here
type Indexer struct {
	sourceFile *os.File
	indexFile  *os.File
	IoLock     sync.Mutex
}

func (self *Indexer) init(indexPath, sourcePath string) *Indexer {
	file, err := os.Open(indexPath)
	if err != nil {
		panic(err.Error())
	}
	self.indexFile = file
	file, err = os.Open(sourcePath)
	if err != nil {
		panic(err.Error())
	}
	self.sourceFile = file
	return self
}

func NewIndexer(indexPath, sourcePath string) *Indexer {
	return new(Indexer).init(indexPath, sourcePath)
}

// GetLine
func (self *Indexer) GetLine(line int64) (string, error) {
	payload, err := self.GetLineBytes(line)
	return string(payload), err
}

// GetLineBytes
func (self *Indexer) GetLineBytes(line int64) ([]byte, error) {
	if line < 0 {
		return []byte{}, ouchErr
	}
	self.IoLock.Lock()
	defer self.IoLock.Unlock()
	from := offsetForLine(line-1, self.indexFile)
	to := offsetForLine(line, self.indexFile)
	if from >= to {
		return []byte{}, offsetErr
	}
	payload := make([]byte, to-from)
	self.sourceFile.Seek(from, SEEK_BEGINNING)
	self.sourceFile.Read(payload)
	return payload, nil
}

// BuildIndex builds a an ofset index from a given file
func BuildIndex(filePath string) (*Indexer, error) {
	sourceFile, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer sourceFile.Close()
	indexPath := filePath + ".index"
	idxFile, err := os.Create(indexPath)

	if err != nil {
		return nil, err
	}
	defer idxFile.Close()
	indexWritter := bufio.NewWriter(idxFile)

	index := make(map[int64]int64)
	var line, offsetBytes int64
	line = 1
	scanner := bufio.NewScanner(sourceFile)
	for scanner.Scan() {
		index[line] = offsetBytes
		err := writeInt32(indexWritter, offsetBytes)
		if err != nil {
			return nil, err
		}
		indexWritter.Flush()
		offsetBytes += int64(len(scanner.Bytes())) + SIZE_OF_NEWLINE
		line++
	}
	log.WithField("name", indexPath).Info("Index has built")
	return NewIndexer(indexPath, filePath), scanner.Err()
}

// Private utility functions

// If you end up reading this code - would you like to hear a TCP/IP Joke? :)
// https://twitter.com/KirkBater/status/953673704734683136
func writeInt32(buf io.Writer, n int64) (err error) {
	payload := make([]byte, SIZE_OF_INT64)
	binary.BigEndian.PutUint64(payload, uint64(n))
	_, err = buf.Write(payload)
	return
}

func readInt32(buf io.Reader) (n int64, err error) {
	payload := make([]byte, SIZE_OF_INT64)
	if _, err := io.ReadFull(buf, payload); err != nil {
		return -1, err
	}
	n = int64(binary.BigEndian.Uint64(payload))
	return
}

func offsetForLine(line int64, file *os.File) int64 {
	file.Seek(line*SIZE_OF_INT64, SEEK_BEGINNING)
	n, _ := readInt32(file)
	return n
}
