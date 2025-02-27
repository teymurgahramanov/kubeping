FROM golang:1.21.6 AS build
WORKDIR /app
COPY . .
RUN \
    if ! test -f bin/kubeping-exporter; then \
      CGO_ENABLED=0 GOOS=linux go build -ldflags "-w -s" -o bin/kubeping-exporter; \
    fi
FROM scratch as final
COPY --from=build /app/bin/kubeping-exporter /
ENTRYPOINT [ "/kubeping-exporter" ]