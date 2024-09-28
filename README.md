
# Practical Networking

Добро пожаловать в репозиторий **Practical Networking**, который предоставляет всестороннее руководство по сетевым технологиям через сочетание теории, практических заданий и реальных примеров. Этот репозиторий будет полезен новичкам, которые хотят углубить свои знания в области сетей.

## Содержание

- [Описание](#описание)
- [Структура репозитория](#структура-репозитория)
- [Теоретический материал](#теоретический-материал)
- [Лабораторные работы](#лабораторные-работы)
- [Реальные задачи](#реальные-задачи)
- [Вспомогательные утилиты](#вспомогательные-утилиты)
- [Как внести вклад](#как-внести-вклад)
- [Лицензия](#лицензия)

---

## Описание

Этот репозиторий создан для того, чтобы предоставить полное понимание сетевых технологий на практическом и теоретическом уровне. Здесь вы найдете:
- **Теоретические материалы** — полные объяснения сетевых протоколов, таких как UDP, TCP, и multicast.
- **Лабораторные работы** — задания, которые помогут закрепить теоретические знания на практике.
- **Реальные задачи** — примеры того, как применить полученные знания в реальных условиях.
- **Утилиты** — инструменты для тестирования и анализа сетевых взаимодействий.

---

## Структура репозитория

```bash
practical-networking/
│
├── README.md                      # Описание репозитория
├── LICENSE                        # Лицензия 
├── docs/                          # Документация и теоретический материал
│   ├── udp-basics.md              # Основы UDP
│   ├── tcp-basics.md              # Основы TCP
│   └── multicast.md               # Теория о multicast и его применении
│
├── labs/                          # Лабораторные работы
│   ├── lab1_udp_multicast/        # Лабораторная работа 1: Работа с UDP Multicast
│   │   ├── README.md              # Описание лабораторной работы и задачи
│   │   ├── solution/              # Решение лабораторной работы
│   │   │   ├── main.go            # Исходный код решения
│   │   └── examples/              # Примеры и тесты для проверки работы
│   ├── lab2_tcp_connection/       # Лабораторная работа 2: TCP соединение
│   │   ├── README.md              # Описание лабораторной работы и задачи
│   │   ├── solution/              # Решение лабораторной работы
│   └── <labN>/                    # Следующие лабораторные работы
│
├── real-world-tasks/              # Реальные задачи на базе лабораторных работ
│   ├── task1_load_balancer/       # Реальная задача 1: Создание балансировщика нагрузки
│   │   ├── README.md              # Описание задачи
│   │   ├── solution/              # Решение задачи
│   └── <taskN>/                   # Следующие задачи
│
└── utils/                         # Вспомогательные утилиты и скрипты
    ├── net-sniffer.go             # Пример утилиты для перехвата сетевых пакетов
    └── <scripts>                  # Другие вспомогательные скрипты
```
Данный пример представлен для лучшего понимания структуры репозитория. Однако стоит отметить, что это лишь демонстрационный вариант, и реальная структура может отличаться.


---

## Теоретический материал

Каталог `docs/` содержит необходимые теоретические материалы для изучения сетевых технологий. Каждая тема подробно рассматривает один из сетевых протоколов.

- **[Основы UDP](./docs/udp-basics.md)** — объяснение работы протокола UDP и примеры его использования.
- **[Основы TCP](./docs/tcp-basics.md)** — обзор протокола TCP и примеры использования.
- **[Multicast](./docs/multicast.md)** — объяснение и использование протоколов multicast.

---

## Лабораторные работы

Каталог `labs/` содержит практические задания, которые помогут вам лучше понять и применить полученные знания на практике.

### Доступные лабораторные работы:
- **[Лабораторная работа 1: UDP Multicast](./labs/lab1_udp_multicast/README.md)** — изучение работы с multicast через протокол UDP.
- **[Лабораторная работа 2: Передача файла через TCP](./labs/lab2_tcp_file_transfer/README.md)** — реализация передачи файла через протокол TCP, поддержка нескольких клиентов.
- **[Лабораторная работа 5: SOCKS-прокси](labs/lab5_socks_proxy/README.md)** — реализация SOCKS5-прокси-сервера с поддержкой установления TCP/IP соединения.

Новые лабораторные работы будут добавляться по мере развития репозитория.

---

## Реальные задачи

В каталоге `real-world-tasks/` представлены реальные задачи, созданные на основе лабораторных работ, которые демонстрируют, как применять знания сетевых технологий в практических приложениях.

- **[Задача 1: Создание балансировщика нагрузки](./real-world-tasks/task1_load_balancer/README.md)** — создание простого балансировщика нагрузки, используя знания, полученные в лабораторных работах.

---

## Вспомогательные утилиты

Каталог `utils/` содержит вспомогательные скрипты и утилиты, которые помогут вам в проведении экспериментов и отладки сетевых приложений.

---

## Как внести вклад

Мы приветствуем любые вклады в развитие этого репозитория! Если у вас есть идеи для новых лабораторных работ, задач или улучшений документации, не стесняйтесь открывать pull request или issue.

1. Сделайте форк репозитория
2. Создайте новую ветку для изменений (`git checkout -b feature/my-new-feature`)
3. Внесите изменения и зафиксируйте их (`git commit -am 'Add my new feature'`)
4. Запушьте изменения (`git push origin feature/my-new-feature`)
5. Откройте Pull Request

Пожалуйста, убедитесь, что ваш код соответствует нашим [рекомендациям по внесению изменений](CONTRIBUTING.md).

---

## Лицензия

Этот репозиторий распространяется под лицензией MIT. Подробности можно найти в файле [LICENSE](./LICENSE).

---

