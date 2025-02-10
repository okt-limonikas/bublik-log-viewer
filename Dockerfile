FROM golang:1.22 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/main.go

FROM alpine:latest

WORKDIR /root/

RUN mkdir json

COPY --from=builder /app/app .

EXPOSE 5050

CMD ["./app", "serve", "./json", "--host", "0.0.0.0"]
