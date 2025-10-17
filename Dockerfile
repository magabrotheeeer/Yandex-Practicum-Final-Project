FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o todo .

FROM alpine:3.19

WORKDIR /app

COPY --from=builder /app/todo .

COPY --from=builder /app/.env ./

EXPOSE 7540

CMD ["./todo"]