FROM golang:1.22 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/blv/main.go

FROM alpine:latest

WORKDIR /root/

RUN mkdir -p /root/json

COPY --from=builder /app/app /usr/local/bin/blv

EXPOSE 5050

ENTRYPOINT ["blv"]
CMD ["serve", "/root/json"]
