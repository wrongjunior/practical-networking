package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	port           = ":9999"         // Порт для обмена сообщениями
	interval       = 2 * time.Second // Интервал между отправкой сообщений
	expiryDuration = 5 * time.Second // Время, через которое копия считается исчезнувшей
)

var (
	peers     = make(map[string]time.Time) // Список IP-адресов живых копий
	peersLock sync.Mutex                   // Для синхронизации доступа к списку
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Использование: go run main.go <multicast-address>")
	}

	multicastAddress := os.Args[1]

	fmt.Printf("Запуск программы с multicast-адресом: %s\n", multicastAddress)

	go receive(multicastAddress)
	go send(multicastAddress)

	// Обновление и вывод списка живых копий
	for {
		time.Sleep(1 * time.Second)
		printAlivePeers()
	}
}

// Отправка сообщения о присутствии
func send(multicastAddress string) {
	addr, err := net.ResolveUDPAddr("udp", multicastAddress+port)
	if err != nil {
		log.Fatalf("Ошибка разрешения адреса: %v", err)
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Fatalf("Ошибка при подключении к UDP: %v", err)
	}
	defer conn.Close()

	fmt.Println("Отправка сообщений о присутствии...")

	for {
		_, err := conn.Write([]byte("I'm here"))
		if err != nil {
			log.Printf("Ошибка отправки сообщения: %v", err)
		} else {
			fmt.Println("Сообщение отправлено")
		}
		time.Sleep(interval)
	}
}

// Прием сообщений от других копий
func receive(multicastAddress string) {
	addr, err := net.ResolveUDPAddr("udp", multicastAddress+port)
	if err != nil {
		log.Fatalf("Ошибка разрешения адреса: %v", err)
	}

	conn, err := net.ListenMulticastUDP("udp", nil, addr)
	if err != nil {
		log.Fatalf("Ошибка при подключении к multicast: %v", err)
	}
	defer conn.Close()

	conn.SetReadBuffer(1024)

	fmt.Println("Ожидание сообщений от других копий...")

	buffer := make([]byte, 1024)
	for {
		n, src, err := conn.ReadFromUDP(buffer)
		if err != nil {
			log.Printf("Ошибка при получении сообщения: %v", err)
			continue
		}

		message := strings.TrimSpace(string(buffer[:n]))
		if message == "I'm here" {
			fmt.Printf("Получено сообщение от: %s\n", src.IP.String())
			updatePeer(src.IP.String())
		}
	}
}

// Обновление времени активности копии
func updatePeer(ip string) {
	peersLock.Lock()
	defer peersLock.Unlock()

	peers[ip] = time.Now()
}

// Печать списка живых копий
func printAlivePeers() {
	peersLock.Lock()
	defer peersLock.Unlock()

	now := time.Now()
	changed := false

	// Удаление неактивных копий
	for ip, lastSeen := range peers {
		if now.Sub(lastSeen) > expiryDuration {
			delete(peers, ip)
			changed = true
		}
	}

	if changed || len(peers) > 0 {
		fmt.Println("Текущие живые копии:")
		for ip := range peers {
			fmt.Printf("- %s\n", ip)
		}
		fmt.Println("---------------------------")
	}
}
