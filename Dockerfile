FROM golang:1.25-alpine AS builder

WORKDIR /app

RUN apk add --no-cache gcc musl-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 go build -o shoplist ./cmd/server/main.go

FROM alpine:latest

ARG DB_URL

ENV DB_URL=$DB_URL

RUN apk add --no-cache libc6-compat ca-certificates

WORKDIR /app

COPY --from=builder /app/shoplist /app/shoplist
COPY --from=builder /app/web /app/web
RUN mkdir -p /app/data

EXPOSE 8080

CMD ["/app/shoplist"]