# scheduler_tgbot

## Зависимости
docker

## Подготовка
```
git clone https://github.com/greg-zamza/scheduler_tgbot
cd scheduler_tgbot
```

Создай файл `.env` для переменных окружения, которые передадутся в контейнеры (пример файла в `.env.example`).

Собери образ scheduler_tgbot:
```
cd BotService
docker build -t scheduler_tgbot .
cd ..
```

## Запуск
```
docker compose up
> или docker compose up -d
```
