FROM golang:1.21 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o avito-slug ./cmd/avito-slug

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/avito-slug .
COPY --from=builder /app/config /app/config

EXPOSE 8080

ENV CONFIG_PATH="config/local.yml"

CMD ["./avito-slug"]