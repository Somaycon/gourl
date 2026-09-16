# Substitua golang:1.22-alpine por golang:alpine
FROM golang:alpine AS builder

WORKDIR /app

ENV GOPROXY=https://proxy.golang.org,direct

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Compila o binário na nova estrutura de pastas
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api/main.go

FROM alpine:latest

WORKDIR /app

RUN apk --no-cache add ca-certificates

COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]
