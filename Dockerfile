FROM golang:1.19-alpine as builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server cmd/server/main.go

FROM alpine:3.14

WORKDIR /app

COPY --from=builder /app/server .

EXPOSE 8080

CMD ["./server"]
