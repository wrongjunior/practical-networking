
# Установка TCP соединения

Когда мы работаем с сетевыми приложениями, передача данных между клиентом и сервером требует создания надёжного соединения. Для этого в TCP используется процесс, называемый "трёхстороннее рукопожатие" (Three-Way Handshake). Этот механизм обеспечивает надёжную и корректную передачу данных, гарантируя, что обе стороны готовы к началу общения.

## Что такое трёхстороннее рукопожатие?

TCP соединение между двумя устройствами устанавливается с помощью процесса, который состоит из трёх шагов:

1. **SYN (synchronize)**: Клиент инициирует соединение и отправляет серверу запрос на установку связи, называемый сегментом SYN. Этот сегмент содержит начальный порядковый номер (sequence number), который будет использоваться для отслеживания отправленных данных.
   
2. **SYN-ACK (synchronize-acknowledge)**: Сервер принимает запрос клиента и отправляет ответный сегмент SYN-ACK. Этот сегмент подтверждает получение SYN от клиента (ACK) и отправляет свой собственный SYN для инициирования связи с клиентом.

3. **ACK (acknowledge)**: Клиент, получив SYN-ACK от сервера, отправляет последний сегмент ACK, подтверждая установку соединения. Теперь оба устройства готовы к передаче данных.

### Пример трёхстороннего рукопожатия:

1. Клиент (IP 192.168.0.1) хочет подключиться к серверу (IP 192.168.0.2) на порту 80. Клиент отправляет SYN-пакет:

    ```
    Клиент -> Сервер: SYN, Seq=100 (я хочу начать передачу)
    ```

2. Сервер принимает этот пакет и отвечает SYN-ACK:

    ```
    Сервер -> Клиент: SYN, Seq=300, ACK=101 (я готов, подтверждаю твой запрос)
    ```

3. Клиент отправляет последний пакет ACK для завершения соединения:

    ```
    Клиент -> Сервер: ACK=301 (я подтверждаю, соединение установлено)
    ```

Теперь соединение установлено, и начинается передача данных. Оба участника (клиент и сервер) знают начальные порядковые номера, что позволяет отслеживать каждую отправленную и полученную часть данных.

## Практическое значение трёхстороннего рукопожатия

- **Надёжность передачи**: Использование трёх пакетов для установления соединения обеспечивает надёжную передачу данных, так как каждое устройство получает подтверждение готовности второго.
- **Инициализация обмена данными**: Устанавливая порядковые номера для обеих сторон, TCP готовится к надёжному обмену данными, где каждый сегмент данных может быть подтверждён.
- **Защита от ошибок**: TCP предусматривает повторную передачу пакетов при сбоях, а начальная фаза с SYN, SYN-ACK и ACK минимизирует вероятность ошибок уже на этапе установления соединения.

## Закрытие соединения (TCP Teardown)

Закрытие TCP-соединения также требует обмена сообщениями, однако оно может быть немного более сложным, так как каждая сторона может завершить соединение независимо.

1. **FIN (finish)**: Когда одна сторона (например, клиент) хочет закрыть соединение, она отправляет сегмент FIN, что означает завершение передачи данных.

2. **ACK (acknowledge)**: Другая сторона (например, сервер) принимает FIN и отправляет сегмент ACK, подтверждая получение сообщения о завершении.

3. **FIN (finish)**: Сервер, после того как обработает все свои данные, может отправить собственное сообщение FIN, чтобы уведомить клиента о закрытии соединения.

4. **ACK (acknowledge)**: Клиент отправляет завершающий ACK, после чего соединение считается закрытым.

### Пример закрытия соединения:

1. Клиент завершает передачу данных и отправляет сегмент FIN:

    ```
    Клиент -> Сервер: FIN, Seq=400 (я закончил передачу данных)
    ```

2. Сервер принимает FIN и отправляет ACK:

    ```
    Сервер -> Клиент: ACK=401 (подтверждаю завершение передачи)
    ```

3. Когда сервер готов завершить, он отправляет свой сегмент FIN:

    ```
    Сервер -> Клиент: FIN, Seq=700 (я тоже закончил передачу)
    ```

4. Клиент подтверждает завершение с помощью ACK:

    ```
    Клиент -> Сервер: ACK=701 (подтверждаю завершение)
    ```

Теперь соединение завершено с обеих сторон.

## Задержки и потери пакетов

Важно учитывать, что в реальной сети могут возникнуть задержки и потери пакетов. TCP автоматически обрабатывает эти ситуации:
- Если клиент не получает ACK от сервера, он повторно отправляет SYN через некоторое время.
- Если сервер не получает ACK от клиента, он также может повторно отправить SYN-ACK.

Таким образом, TCP защищает соединение от потерь и задержек, обеспечивая надёжность.

## Заключение

Механизм трёхстороннего рукопожатия является основой работы протокола TCP, позволяя устанавливать соединение между двумя устройствами в сети. Он обеспечивает надёжное установление связи, обмен данными и их подтверждение, что делает TCP одним из самых надёжных транспортных протоколов в сетях.
