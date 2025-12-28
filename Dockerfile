FROM golang:1.25.1-alpine AS builder

WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o server .

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/server .
EXPOSE 8080
CMD ["./server"]
