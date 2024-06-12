FROM docker.io/library/golang:1.22.4-alpine3.20 AS builder

WORKDIR /builder
COPY ./src/ /builder/
RUN CGO_ENABLED=0 go build -o dns-updater -ldflags '-extldflags "-static" -w -s'  ./...

FROM docker.io/library/alpine:3.20

COPY --from=builder /builder/dns-updater /bin/dns-updater
USER 1000

ENTRYPOINT [ "/bin/dns-updater"]