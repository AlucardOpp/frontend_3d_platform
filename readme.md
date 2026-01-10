## Запуск

Обновить образы контейнеров
```bash
docker-compose -f docker-compose.yml pull
```

Для запуска на сервере
```bash
docker-compose up -d
```

Для остановки на сервере
```bash
docker-compose down
```

Запуск для локальной разработки (с пробросом портов)
```bash
docker-compose -f docker-compose.yml -f docker-compose.local.yml up -d --build
```

Остановка для локальной разработки
```bash
docker-compose -f docker-compose.yml -f docker-compose.local.yml down
```