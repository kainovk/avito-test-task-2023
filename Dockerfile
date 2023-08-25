FROM golang:1.21
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY config/local.yml ./config/local.yml
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o avito-slug ./cmd/avito-slug
EXPOSE 8080
CMD ["./avito-slug"]