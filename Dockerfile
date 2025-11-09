FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o bookvault-api main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/bookvault-api .

COPY .env .

EXPOSE 8000

CMD ["./bookvault-api"]
