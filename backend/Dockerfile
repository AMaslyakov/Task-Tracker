FROM golang:1.25-alpine AS builder

WORKDIR /src

COPY src/go.mod src/go.sum ./
RUN go mod download
RUN go install github.com/swaggo/swag/cmd/swag@v1.16.6

COPY src/ ./
RUN swag init -g main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/backend .

FROM alpine:3.22

WORKDIR /app

COPY --from=builder /app/backend ./backend

EXPOSE 8080

CMD ["./backend"]
