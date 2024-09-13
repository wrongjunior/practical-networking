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
	port           = ":9999"          // порт для обмена сообщениями
	interval       = 4 * time.Second  // интервал между отправкой сообщений
	expiryDuration = 10 * time.Second // время, через которое копия считается исчезнувшей
)

var (
	peers     = make(map[string]time.Time) // список IP-адресов живых копий
	peersLock sync.Mutex                   // для синхронизации доступа к списку
	myID      = uuid.New().String()        // уникальный идентификатор текущей копии
)

// Config хранит конфигурацию приложения
type Config struct {
	Protocol         string // протокол (IPv4 или IPv6)
	MulticastAddress string // адрес multicast-группы
	InterfaceName    string // имя сетевого интерфейса (для IPv6)
}

func main() {
	// разбираем аргументы командной строки и получаем конфигурацию
	config, err := parseArgs()
	if err != nil {
		log.Fatalf("Ошибка при разборе аргументов: %v", err)
	}

	// получаем сетевой интерфейс для использования
	iface, err := getInterface(config)
	if err != nil {
		log.Fatalf("Ошибка при получении интерфейса: %v", err)
	}

	// определяем multicast-адрес
	addr, err := net.ResolveUDPAddr("udp", config.MulticastAddress+port)
	if err != nil {
		log.Fatalf("Ошибка при разрешении адреса: %v", err)
	}

	fmt.Printf("Запуск приложения с %s multicast-адресом: %s, интерфейс: %s, ID копии: %s\n",
		config.Protocol, config.MulticastAddress, config.InterfaceName, myID)

	// запускаем прием и отправку сообщений в отдельных горутинах
	go receive(addr, iface)
	go send(addr, iface)

	// обновление и вывод списка живых копий
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
		MulticastAddress: os.Args[1], // получаем multicast-адрес
	}

	// определяем протокол в зависимости от наличия символа ":" в адресе
	if strings.Contains(config.MulticastAddress, ":") {
		config.Protocol = "ipv6"
		config.MulticastAddress = "[" + config.MulticastAddress + "]" // для IPv6 адрес обрамляем в скобки
		if len(os.Args) < 3 {
			return nil, fmt.Errorf("для IPv6 необходимо указать имя интерфейса")
		}
		config.InterfaceName = os.Args[2] // получаем имя интерфейса
	} else {
		config.Protocol = "ipv4"
	}

	return config, nil
}

// getInterface возвращает сетевой интерфейс для использования
func getInterface(config *Config) (*net.Interface, error) {
	if config.Protocol == "ipv6" {
		return net.InterfaceByName(config.InterfaceName) // для IPv6 требуется интерфейс
	}
	return nil, nil // для IPv4 интерфейс не требуется
}

// send отправляет сообщения о присутствии
func send(addr *net.UDPAddr, iface *net.Interface) {
	// подключаемся к UDP
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Fatalf("Ошибка при подключении к UDP: %v", err)
	}
	defer conn.Close()

	fmt.Println("Отправка сообщений о присутствии...")

	// отправляем сообщения в цикле
	for {
		_, err := conn.Write([]byte(fmt.Sprintf("Я здесь, ID: %s", myID))) // отправляем сообщение с идентификатором
		if err != nil {
			log.Printf("Ошибка отправки сообщения: %v", err)
		}
		time.Sleep(interval) // ждем перед следующей отправкой
	}
}

// receive принимает сообщения от других копий
func receive(addr *net.UDPAddr, iface *net.Interface) {
	// слушаем multicast-адрес
	conn, err := net.ListenMulticastUDP("udp", iface, addr)
	if err != nil {
		log.Fatalf("Ошибка при подключении к multicast: %v", err)
	}
	defer conn.Close()

	conn.SetReadBuffer(1024) // устанавливаем размер буфера

	fmt.Println("Ожидание сообщений от других копий...")

	buffer := make([]byte, 1024)
	for {
		n, src, err := conn.ReadFromUDP(buffer) // читаем сообщения
		if err != nil {
			log.Printf("Ошибка при получении сообщения: %v", err)
			continue
		}

		// обрабатываем сообщение
		message := strings.TrimSpace(string(buffer[:n]))
		if strings.Contains(message, "Я здесь") && !strings.Contains(message, myID) {
			fmt.Printf("Получено сообщение от: %s (%s)\n", src.IP.String(), message)
			updatePeer(src.IP.String(), message) // обновляем информацию о копии
		}
	}
}

// updatePeer обновляет время последнего обнаружения для копии
func updatePeer(ip, id string) {
	peersLock.Lock()
	defer peersLock.Unlock()

	peers[ip+" ("+id+")"] = time.Now() // обновляем время для данного IP
}

// checkAlivePeers проверяет и выводит список активных копий
func checkAlivePeers() {
	peersLock.Lock()
	defer peersLock.Unlock()

	now := time.Now()
	changed := false
	var activePeers []string

	// проверяем время последнего обнаружения для каждой копии
	for ip, lastSeen := range peers {
		if now.Sub(lastSeen) > expiryDuration {
			fmt.Printf("Копия %s больше не активна, удаляем...\n", ip)
			delete(peers, ip) // удаляем копию, если она больше не активна
			changed = true
		} else {
			activePeers = append(activePeers, ip)
		}
	}

	// выводим список активных копий, если есть изменения
	if changed || len(activePeers) > 0 {
		fmt.Println("Текущие активные копии:")
		for _, peer := range activePeers {
			fmt.Printf("- %s (последнее обнаружение: %s назад)\n", peer, now.Sub(peers[peer]).Round(time.Second))
		}
		fmt.Println("---------------------------")
	}
}
