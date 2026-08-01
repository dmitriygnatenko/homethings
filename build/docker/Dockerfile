FROM golang:alpine AS builder

WORKDIR ./../../

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o /app cmd/app/main.go

FROM alpine:3.24

WORKDIR /

COPY --from=builder /app /app

EXPOSE 8080

CMD ["/app"]