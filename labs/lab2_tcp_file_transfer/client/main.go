// client.go
package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
)

func main() {
	filePath := flag.String("file", "", "Path to file to send")
	serverHost := flag.String("host", "", "Server host")
	serverPort := flag.String("port", "", "Server port")
	flag.Parse()

	if *filePath == "" || *serverHost == "" || *serverPort == "" {
		log.Fatal("Please provide -file, -host, and -port parameters")
	}

	file, err := os.Open(*filePath)
	if err != nil {
		log.Fatalf("Failed to open file %s: %v", *filePath, err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		log.Fatalf("Failed to get file info: %v", err)
	}

	fileSize := fileInfo.Size()
	if fileSize > 1<<40 {
		log.Fatalf("File size exceeds 1 TB")
	}

	filename := filepath.Base(*filePath)
	filenameBytes := []byte(filename)
	if len(filenameBytes) > 4096 {
		log.Fatalf("Filename exceeds 4096 bytes")
	}

	conn, err := net.Dial("tcp", net.JoinHostPort(*serverHost, *serverPort))
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	nameLen := uint32(len(filenameBytes))
	err = binary.Write(conn, binary.BigEndian, nameLen)
	if err != nil {
		log.Fatalf("Failed to send filename length: %v", err)
	}

	n, err := conn.Write(filenameBytes)
	if err != nil || n != len(filenameBytes) {
		log.Fatalf("Failed to send filename: %v", err)
	}

	err = binary.Write(conn, binary.BigEndian, uint64(fileSize))
	if err != nil {
		log.Fatalf("Failed to send file size: %v", err)
	}

	buf := make([]byte, 32*1024)
	for {
		n, err := file.Read(buf)
		if err != nil && err != io.EOF {
			log.Fatalf("Failed to read from file: %v", err)
		}
		if n == 0 {
			break
		}
		_, err = conn.Write(buf[:n])
		if err != nil {
			log.Fatalf("Failed to send file data: %v", err)
		}
	}

	resp := make([]byte, 1)
	n, err = conn.Read(resp)
	if err != nil {
		log.Fatalf("Failed to receive response from server: %v", err)
	}
	if n != 1 {
		log.Fatalf("Invalid response from server")
	}

	if resp[0] == 0x00 {
		fmt.Println("File transfer successful")
	} else {
		fmt.Println("File transfer failed")
	}
}
