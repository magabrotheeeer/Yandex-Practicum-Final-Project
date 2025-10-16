# Используем официальный Go-образ
FROM golang:1.24-alpine

WORKDIR /app

# Копируем файлы модуля и загружаем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь проект (включая .env) и собираем
COPY . .
RUN go build -o todo .

# Открываем порт
EXPOSE 7540

# Запускаем приложение (оно само прочитает .env)
CMD ["./todo"]