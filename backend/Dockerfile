FROM golang:1.25-alpine AS builder

WORKDIR /src

COPY src/go.mod src/go.sum ./
RUN go mod download

COPY src/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/backend .

FROM alpine:3.22

WORKDIR /app

COPY --from=builder /app/backend ./backend

EXPOSE 8080

CMD ["./backend"]
