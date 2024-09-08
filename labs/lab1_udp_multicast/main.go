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
	port           = ":9999"          // Порт для обмена сообщениями
	interval       = 4 * time.Second  // Интервал между отправкой сообщений
	expiryDuration = 10 * time.Second // Время, через которое копия считается исчезнувшей
)

var (
	peers     = make(map[string]time.Time) // Карта для хранения информации о других копиях
	peersLock sync.Mutex                   // Мьютекс для синхронизации доступа к карте peers
	myID      = uuid.New().String()        // Уникальный идентификатор текущей копии
)

// Config хранит конфигурацию приложения
type Config struct {
	Protocol         string // Протокол (IPv4 или IPv6)
	MulticastAddress string // Адрес multicast-группы
	InterfaceName    string // Имя сетевого интерфейса (для IPv6)
}

func main() {
	config, err := parseArgs()
	if err != nil {
		log.Fatalf("Ошибка при разборе аргументов: %v", err)
	}

	iface, err := getInterface(config)
	if err != nil {
		log.Fatalf("Ошибка при получении интерфейса: %v", err)
	}

	addr, err := net.ResolveUDPAddr("udp", config.MulticastAddress+port)
	if err != nil {
		log.Fatalf("Ошибка при разрешении адреса: %v", err)
	}

	fmt.Printf("Запуск приложения с %s multicast-адресом: %s, интерфейс: %s, ID копии: %s\n",
		config.Protocol, config.MulticastAddress, config.InterfaceName, myID)

	go receive(addr, iface)
	go send(addr, iface)

	// Обновление и вывод списка живых копий
	for {
		time.Sleep(1 * time.Second)
		checkAlivePeers()
	}
}

// parseArgs разбирает аргументы командной строки и возвращает Config
func parseArgs() (*Config, error) {
	if len(os.Args) < 2 {
		return nil, fmt.Errorf("использование: %s <multicast-адрес> [имя-интерфейса]", os.Args[0])
	}

	config := &Config{
		MulticastAddress: os.Args[1],
	}

	if strings.Contains(config.MulticastAddress, ":") {
		config.Protocol = "ipv6"
		config.MulticastAddress = "[" + config.MulticastAddress + "]"
		if len(os.Args) < 3 {
			return nil, fmt.Errorf("для IPv6 необходимо указать имя интерфейса")
		}
		config.InterfaceName = os.Args[2]
	} else {
		config.Protocol = "ipv4"
	}

	return config, nil
}

// getInterface возвращает сетевой интерфейс для использования
func getInterface(config *Config) (*net.Interface, error) {
	if config.Protocol == "ipv6" {
		return net.InterfaceByName(config.InterfaceName)
	}
	return nil, nil
}

// send отправляет сообщения о присутствии
func send(addr *net.UDPAddr, iface *net.Interface) {
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Fatalf("Ошибка при подключении к UDP: %v", err)
	}
	defer conn.Close()

	fmt.Println("Отправка сообщений о присутствии...")

	for {
		_, err := conn.Write([]byte(fmt.Sprintf("Я здесь, ID: %s", myID)))
		if err != nil {
			log.Printf("Ошибка отправки сообщения: %v", err)
		}
		time.Sleep(interval)
	}
}

// receive принимает сообщения от других копий
func receive(addr *net.UDPAddr, iface *net.Interface) {
	conn, err := net.ListenMulticastUDP("udp", iface, addr)
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
		if strings.Contains(message, "Я здесь") && !strings.Contains(message, myID) {
			fmt.Printf("Получено сообщение от: %s (%s)\n", src.IP.String(), message)
			updatePeer(src.IP.String(), message)
		}
	}
}

// updatePeer обновляет время последнего обнаружения для копии
func updatePeer(ip, id string) {
	peersLock.Lock()
	defer peersLock.Unlock()

	peers[ip+" ("+id+")"] = time.Now()
}

// checkAlivePeers проверяет и выводит список активных копий
func checkAlivePeers() {
	peersLock.Lock()
	defer peersLock.Unlock()

	now := time.Now()
	changed := false
	var activePeers []string

	for ip, lastSeen := range peers {
		if now.Sub(lastSeen) > expiryDuration {
			fmt.Printf("Копия %s больше не активна, удаляем...\n", ip)
			delete(peers, ip)
			changed = true
		} else {
			activePeers = append(activePeers, ip)
		}
	}

	if changed || len(activePeers) > 0 {
		fmt.Println("Текущие активные копии:")
		for _, peer := range activePeers {
			fmt.Printf("- %s (последнее обнаружение: %s назад)\n", peer, now.Sub(peers[peer]).Round(time.Second))
		}
		fmt.Println("---------------------------")
	}
}
