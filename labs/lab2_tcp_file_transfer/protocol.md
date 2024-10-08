
# Формализованное описание протокола передачи файла по TCP

## Шаги протокола

### 1. Установление TCP-соединения

Клиент:
- Открывает сокет и подключается к адресу сервера:
    ```go
    conn = net.Dial("tcp", "<адрес_сервера>:<порт>")
    ```

### 2. Передача длины имени файла

Клиент:
- Отправляет серверу длину имени файла (4 байта, целое число в формате Big Endian):
    ```go
    nameLen = uint32(len(filename))
    binary.Write(conn, binary.BigEndian, nameLen)
    ```

Сервер:
- Принимает 4 байта и интерпретирует их как длину имени файла:
    ```go
    binary.Read(conn, binary.BigEndian, &nameLen)
    ```

### 3. Передача имени файла

Клиент:
- Отправляет само имя файла (строка длиной `nameLen` байт):
    ```go
    conn.Write([]byte(filename))
    ```

Сервер:
- Читает `nameLen` байт и интерпретирует их как имя файла:
    ```go
    io.ReadFull(conn, nameBytes)
    ```

### 4. Передача размера файла

Клиент:
- Отправляет серверу размер файла (8 байт, целое число в формате Big Endian):
    ```go
    fileSize = fileInfo.Size()
    binary.Write(conn, binary.BigEndian, fileSize)
    ```

Сервер:
- Принимает 8 байт и интерпретирует их как размер файла:
    ```go
    binary.Read(conn, binary.BigEndian, &fileSize)
    ```

### 5. Передача содержимого файла

Клиент:
- Отправляет файл частями по 32 КБ:
    ```go
    for {
        n, _ := file.Read(buf)
        if n == 0 {
            break
        }
        conn.Write(buf[:n])
    }
    ```

Сервер:
- Читает данные и записывает их в файл, пока не получит все байты:
    ```go
    for bytesReceived < fileSize {
        conn.Read(buf)
        file.Write(buf)
    }
    ```

### 6. Отправка сигнала завершения

Сервер:
- По завершению передачи отправляет клиенту сигнал об успешной передаче (1 байт: `0x00` — успех, `0x01` — ошибка):
    ```go
    if bytesReceived == fileSize {
        conn.Write([]byte{0x00})
    } else {
        conn.Write([]byte{0x01})
    }
    ```

Клиент:
- Принимает сигнал и проверяет результат:
    ```go
    conn.Read(resp)
    if resp[0] == 0x00 {
        fmt.Println("File transfer successful")
    } else {
        fmt.Println("File transfer failed")
    }
    ```

## Общие принципы

- **TCP** обеспечивает надёжную передачу данных, включая гарантированную доставку и сохранение порядка передачи.
- Протокол поверх TCP управляет:
  - **Метаданными**: передача длины имени файла, самого имени файла и его размера.
  - **Передачей данных**: отправка файла частями.
  - **Подтверждением**: сигнал об успешной или неуспешной передаче файла.

