FROM golang:1.23 AS builder
WORKDIR /app
COPY . .
RUN go build -o paxos-server web_service.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/paxos-server .
CMD ["./paxos-server"]
