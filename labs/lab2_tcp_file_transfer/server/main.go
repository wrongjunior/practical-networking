// server.go
package main

import (
	"encoding/binary"
	"flag"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

func main() {
	port := flag.String("port", "", "Port to listen on")
	flag.Parse()

	if *port == "" {
		log.Fatal("Please provide a port number using -port")
	}

	// Create uploads directory if it doesn't exist
	err := os.MkdirAll("uploads", os.ModePerm)
	if err != nil {
		log.Fatalf("Failed to create uploads directory: %v", err)
	}

	ln, err := net.Listen("tcp", ":"+*port)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", *port, err)
	}
	defer ln.Close()

	log.Printf("Server listening on port %s", *port)

	var wg sync.WaitGroup
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}
		wg.Add(1)
		go handleConnection(conn, &wg)
	}
}

func handleConnection(conn net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()
	defer conn.Close()

	clientAddr := conn.RemoteAddr().String()
	log.Printf("Connection accepted from %s", clientAddr)

	startTime := time.Now()
	done := make(chan bool)
	var totalBytes uint64 = 0
	var intervalBytes uint64 = 0
	var mutex sync.Mutex

	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	go func() {
		for {
			select {
			case <-ticker.C:
				mutex.Lock()
				elapsed := time.Since(startTime).Seconds()
				avgSpeed := float64(totalBytes) / elapsed
				instSpeed := float64(intervalBytes) / 3.0
				log.Printf("Client %s - Instantaneous speed: %.2f bytes/s, Average speed: %.2f bytes/s", clientAddr, instSpeed, avgSpeed)
				intervalBytes = 0
				mutex.Unlock()
			case <-done:
				mutex.Lock()
				elapsed := time.Since(startTime).Seconds()
				if elapsed == 0 {
					elapsed = 0.001
				}
				avgSpeed := float64(totalBytes) / elapsed
				instSpeed := float64(intervalBytes) / elapsed
				log.Printf("Client %s - Instantaneous speed: %.2f bytes/s, Average speed: %.2f bytes/s", clientAddr, instSpeed, avgSpeed)
				mutex.Unlock()
				return
			}
		}
	}()

	var nameLen uint32
	err := binary.Read(conn, binary.BigEndian, &nameLen)
	if err != nil {
		log.Printf("Failed to read filename length from %s: %v", clientAddr, err)
		return
	}
	totalBytes += 4
	intervalBytes += 4

	if nameLen > 4096 {
		log.Printf("Filename length too long from %s", clientAddr)
		return
	}

	nameBytes := make([]byte, nameLen)
	n, err := io.ReadFull(conn, nameBytes)
	if err != nil {
		log.Printf("Failed to read filename from %s: %v", clientAddr, err)
		return
	}
	totalBytes += uint64(n)
	intervalBytes += uint64(n)

	filename := string(nameBytes)
	filename = filepath.Base(filename)
	filename = strings.TrimSpace(filename)
	if filename == "" {
		log.Printf("Invalid filename from %s", clientAddr)
		return
	}

	var fileSize uint64
	err = binary.Read(conn, binary.BigEndian, &fileSize)
	if err != nil {
		log.Printf("Failed to read file size from %s: %v", clientAddr, err)
		return
	}
	totalBytes += 8
	intervalBytes += 8

	if fileSize > 1<<40 {
		log.Printf("File size exceeds 1 TB from %s", clientAddr)
		return
	}

	uploadPath := filepath.Join("uploads", filename)
	file, err := os.Create(uploadPath)
	if err != nil {
		log.Printf("Failed to create file %s: %v", uploadPath, err)
		return
	}
	defer file.Close()

	buf := make([]byte, 32*1024)
	var bytesReceived uint64 = 0
	for bytesReceived < fileSize {
		toRead := len(buf)
		if remaining := fileSize - bytesReceived; uint64(toRead) > remaining {
			toRead = int(remaining)
		}
		n, err := conn.Read(buf[:toRead])
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Printf("Error reading from %s: %v", clientAddr, err)
			return
		}
		if n > 0 {
			_, err := file.Write(buf[:n])
			if err != nil {
				log.Printf("Error writing to file %s: %v", uploadPath, err)
				return
			}
			mutex.Lock()
			totalBytes += uint64(n)
			intervalBytes += uint64(n)
			bytesReceived += uint64(n)
			mutex.Unlock()
		}
	}

	done <- true

	if bytesReceived != fileSize {
		log.Printf("Received file size does not match expected size from %s", clientAddr)
		conn.Write([]byte{0x01})
		return
	}

	conn.Write([]byte{0x00})
	log.Printf("Successfully received file %s from %s", filename, clientAddr)
}
