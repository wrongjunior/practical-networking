package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

const (
	port           = ":9999"         // Порт для обмена сообщениями
	interval       = 2 * time.Second // Интервал между отправкой сообщений
	expiryDuration = 5 * time.Second // Время, через которое копия считается исчезнувшей
)

var (
	peers     = make(map[string]time.Time) // Список IP-адресов живых копий
	peersLock sync.Mutex                   // Для синхронизации доступа к списку
	myID      = uuid.New().String()        // Уникальный идентификатор текущей копии
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Использование: go run main.go <multicast-address>")
	}

	multicastAddress := os.Args[1]
	fmt.Printf("Запуск программы с multicast-адресом: %s, ID копии: %s\n", multicastAddress, myID)

	go receive(multicastAddress)
	go send(multicastAddress)

	// Обновление и вывод списка живых копий
	for {
		time.Sleep(1 * time.Second)
		checkAlivePeers()
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
		_, err := conn.Write([]byte(fmt.Sprintf("I'm here, ID: %s", myID)))
		if err != nil {
			log.Printf("Ошибка отправки сообщения: %v", err)
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
		if strings.Contains(message, "I'm here") && !strings.Contains(message, myID) {
			fmt.Printf("Получено сообщение от: %s (%s)\n", src.IP.String(), message)
			updatePeer(src.IP.String(), message)
		}
	}
}

// Обновление времени активности копии
func updatePeer(ip, id string) {
	peersLock.Lock()
	defer peersLock.Unlock()

	peers[ip+" ("+id+")"] = time.Now()
}

// Проверка и вывод списка живых копий
func checkAlivePeers() {
	peersLock.Lock()
	defer peersLock.Unlock()

	now := time.Now()
	changed := false
	var activePeers []string

	// Удаление неактивных копий
	for ip, lastSeen := range peers {
		if now.Sub(lastSeen) > expiryDuration {
			fmt.Printf("Копия %s более неактивна, удаление...\n", ip)
			delete(peers, ip)
			changed = true
		} else {
			activePeers = append(activePeers, ip)
		}
	}

	// Вывод активных копий, если произошли изменения
	if changed || len(activePeers) > 0 {
		fmt.Println("Текущие активные копии:")
		for _, peer := range activePeers {
			fmt.Printf("- %s (обнаружено: %s назад)\n", peer, now.Sub(peers[peer]).Round(time.Second))
		}
		fmt.Println("---------------------------")
	}
}
