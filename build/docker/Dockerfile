FROM alpine:3.20

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY build/app/app /app/app

RUN chmod +x /app/app

EXPOSE 8080

ENTRYPOINT ["/app/app"]