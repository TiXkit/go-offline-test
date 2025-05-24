FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download #кэширование

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/main ./app/app.go

FROM scratch

WORKDIR /app

COPY --from=builder /app/main /app/main

ENTRYPOINT ["/app/main"]