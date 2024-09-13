package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
)

func main() {
	numClients := flag.Int("clients", 5, "Number of simultaneous clients")
	serverHost := flag.String("host", "localhost", "Server host")
	serverPort := flag.String("port", "8080", "Server port")
	flag.Parse()

	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "file-transfer-test")
	if err != nil {
		log.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	var wg sync.WaitGroup
	for i := 0; i < *numClients; i++ {
		wg.Add(1)
		go func(clientID int) {
			defer wg.Done()

			// Create a random file
			fileName := fmt.Sprintf("test_file_%d.txt", clientID)
			filePath := filepath.Join(tempDir, fileName)
			fileSize := rand.Intn(1024*1024) + 1024 // Random size between 1KB and 1MB
			err := createRandomFile(filePath, fileSize)
			if err != nil {
				log.Printf("Client %d: Failed to create test file: %v", clientID, err)
				return
			}

			// Run the client
			cmd := exec.Command("go", "run", "/Users/daniilsolovey/Program/go/practical-networking/labs/lab2_tcp_file_transfer/client/main.go",
				"-file", filePath,
				"-host", *serverHost,
				"-port", *serverPort)

			output, err := cmd.CombinedOutput()
			if err != nil {
				log.Printf("Client %d: Failed to run client: %v", clientID, err)
				log.Printf("Client %d output: %s", clientID, string(output))
			} else {
				log.Printf("Client %d: Successfully transferred file. Output: %s", clientID, string(output))
			}
		}(i)
	}

	wg.Wait()
	fmt.Println("All clients finished")
}

func createRandomFile(filePath string, size int) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	data := make([]byte, size)
	rand.Read(data)
	_, err = file.Write(data)
	return err
}
