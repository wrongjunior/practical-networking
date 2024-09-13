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
	// получаем количество клиентов, хост сервера и порт из аргументов командной строки
	numClients := flag.Int("clients", 5, "Number of simultaneous clients")
	serverHost := flag.String("host", "localhost", "Server host")
	serverPort := flag.String("port", "8080", "Server port")
	flag.Parse()

	// создаем временную директорию для тестовых файлов
	tempDir, err := os.MkdirTemp("", "file-transfer-test")
	if err != nil {
		log.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir) // удаляем временную директорию после завершения работы

	var wg sync.WaitGroup
	for i := 0; i < *numClients; i++ {
		wg.Add(1)
		go func(clientID int) {
			defer wg.Done()

			// создаем случайный файл для каждого клиента
			fileName := fmt.Sprintf("test_file_%d.txt", clientID)
			filePath := filepath.Join(tempDir, fileName)
			fileSize := rand.Intn(1024*1024) + 1024 // случайный размер файла от 1КБ до 1МБ
			err := createRandomFile(filePath, fileSize)
			if err != nil {
				log.Printf("Client %d: Failed to create test file: %v", clientID, err)
				return
			}

			// запускаем клиент для передачи файла
			cmd := exec.Command("go", "run", "/Users/daniilsolovey/Program/go/practical-networking/labs/lab2_tcp_file_transfer/client/main.go",
				"-file", filePath,
				"-host", *serverHost,
				"-port", *serverPort)

			output, err := cmd.CombinedOutput() // собираем вывод клиента
			if err != nil {
				log.Printf("Client %d: Failed to run client: %v", clientID, err)
				log.Printf("Client %d output: %s", clientID, string(output)) // вывод ошибки, если есть
			} else {
				log.Printf("Client %d: Successfully transferred file. Output: %s", clientID, string(output)) // успешный вывод
			}
		}(i)
	}

	wg.Wait() // ждем завершения всех клиентов
	fmt.Println("All clients finished")
}

// функция для создания случайного файла заданного размера
func createRandomFile(filePath string, size int) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	data := make([]byte, size)
	rand.Read(data) // заполняем файл случайными данными
	_, err = file.Write(data)
	return err
}
