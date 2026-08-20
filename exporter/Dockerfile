FROM golang:1.27.0 AS build
WORKDIR /app
COPY . .
RUN \
    if ! test -f bin/kubeping-exporter; then \
      CGO_ENABLED=0 GOOS=linux go build -ldflags "-w -s" -o bin/kubeping-exporter; \
    fi

FROM alpine:3.19 AS certs
RUN apk --no-cache add ca-certificates

FROM scratch AS final
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
ENV SSL_CERT_FILE=/etc/ssl/certs/ca-certificates.crt
COPY --from=build /app/bin/kubeping-exporter /
ENTRYPOINT [ "/kubeping-exporter" ]
