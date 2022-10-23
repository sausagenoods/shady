FROM golang:1.19-alpine3.16 as builder
WORKDIR /build
COPY . .
RUN CGO_ENABLED=0 go build -ldflags "-s -w" -o shady-server .

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
WORKDIR /app
COPY --from=builder /build/shady-server .
ENTRYPOINT ["./shady-server", "-bind=:1337"]
