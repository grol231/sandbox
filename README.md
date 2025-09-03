# StarLine RabbitMQ Worker

Go воркер для чтения сообщений из RabbitMQ и отправки запросов в StarLine API.

## Структура проекта

Проект следует [golang-standards/project-layout](https://github.com/golang-standards/project-layout):

```
├── cmd/worker/           # Точка входа приложения
├── internal/             # Внутренние пакеты
│   ├── api/             # HTTP клиент для API
│   ├── config/          # Управление конфигурацией
│   ├── logging/         # Система логирования
│   ├── metrics/         # Prometheus метрики
│   └── worker/          # Основная логика воркера
├── configs/             # Конфигурационные файлы
├── build/               # Dockerfile и скрипты сборки
└── .vscode/             # Конфигурация для VS Code

```

## Конфигурация

Приложение использует YAML конфигурацию. Пример файла `configs/config.yaml`:

```yaml
rabbitmq:
  host: rmq18-prod-slon.sl.netlo
  port: 5672
  user: shinkevich
  password: RANDOM_STRING
  queue: sms

api:
  url: https://lk.zagruzka.com/Starline_http
  service_id: Starline_http
  pass: RANDOM_STRING
  source: StarLine

server:
  port: 8080
  metrics_path: /metrics

logging:
  level: info
  format: json
```

## Сборка и запуск

### Локальная сборка

```bash
# Скачать зависимости
make deps

# Собрать приложение
make build

# Запустить локально
make run

# Запустить тесты
make test

# Запустить тесты с покрытием
make test-coverage
```

### Docker

```bash
# Собрать Docker образ
make docker-build

# Запустить в контейнере
make docker-run
```

## Отладка

Проект включает конфигурацию для VS Code. Используйте:
- `Launch StarLine Worker` - для обычного запуска
- `Debug StarLine Worker` - для отладки

## Мониторинг

### Prometheus метрики

Приложение экспортирует метрики на порту 8080:

- `rabbitmq_messages_received_total` - количество полученных сообщений
- `messages_processed_total` - количество обработанных сообщений  
- `api_requests_sent_total` - количество отправленных API запросов
- `api_requests_success_total` - количество успешных API запросов
- `api_requests_failed_total` - количество неудачных API запросов
- `message_processing_duration_seconds` - время обработки сообщений
- `api_request_duration_seconds` - время выполнения API запросов
- `worker_healthy` - статус здоровья воркера (1 = здоров, 0 = нездоров)

### Логирование

Логи выводятся в JSON формате для удобного парсинга Loki. Все логи на английском языке.

### Health Check

Эндпоинт `/health` возвращает статус приложения.

## Формат сообщений

RabbitMQ сообщения должны быть в JSON формате:

```json
{
  "messages": [
    {
      "recipient": "79218897127",
      "body": "StarLine код авторизации: 2652"
    }
  ]
}
```

Где:
- `recipient` - clientId для API запроса
- `body` - текст сообщения

## API

Воркер отправляет POST запросы на `https://lk.zagruzka.com/Starline_http` со следующими параметрами:
- `clientId` - из поля recipient
- `message` - из поля body  
- `serviceId=Starline_http`
- `pass=RANDOM_STRING` (из конфига)
- `source=StarLine`