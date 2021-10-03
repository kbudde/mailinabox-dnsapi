FROM alpine AS builder
RUN apk add --no-cache ca-certificates

FROM scratch AS final
LABEL maintainer="Kris@budd.ee"
COPY  --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY ./mailinabox-dnsapi /
EXPOSE 8080
USER 65535:65535
ENTRYPOINT ["/mailinabox-dnsapi"]