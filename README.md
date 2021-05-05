# ✉️ mailinabox-dnsapi

This tool allows creating [Let's Encrypt](https://letsencrypt.org/) certificates in [Traefik](https://doc.traefik.io/traefik/https/acme/) if you are using the dns server from [Mail-in-a-Box](https://mailinabox.email/) for managing your domain.

All you have to do is to point traffic with a [httpreq (lego)](https://go-acme.github.io/lego/dns/httpreq/) to this service. This service will convert the request and create a DNS txt record for your domain to proove you ownership.

You can find an example deployment file in [k8s/mailinabox.yaml](k8s/mailinabox.yaml)

## 📝 Todos

- [ ] Create example traefik configuration

