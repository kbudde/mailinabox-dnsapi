# ✉️ mailinabox-dnsapi

This tool allows creating [Let's Encrypt](https://letsencrypt.org/) certificates with [Traefik](https://doc.traefik.io/traefik/https/acme/) if you are using the dns server from [Mail-in-a-Box](https://mailinabox.email/) for managing your domain.

All you have to do is to point configure traefik with an acme resolver with [httpreq (lego)](https://go-acme.github.io/lego/dns/httpreq/) to this service. This service will convert the request and create a DNS txt record for your domain to prove you ownership.

You can find an example deployment file in [k8s/mailinabox.yaml](k8s/mailinabox.yaml)


## traefik configuration snippets

### traefik deployment

```
...
containers:
- image: traefik:v2.4
  name: traefik
  env:
    - name: HTTPREQ_ENDPOINT
      value: http://mailinabox:8080
```

### traefik config file
```
...
certificatesResolvers:
  myresolver:
    acme:
      email: "you@yourdomain.com"
      storage: "acme.json"
      dnsChallenge:
        provider: httpreq
```
### traefik ingress

```
---
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: pihole-https
  namespace: pihole
  annotations:
    traefik.ingress.kubernetes.io/router.entrypoints: websecure
    traefik.ingress.kubernetes.io/router.tls: "true"
    traefik.ingress.kubernetes.io/router.tls.certresolver: myresolver

spec:
  tls:
  - hosts:
      - pihole.yourdomain.com
  rules:
  - host: pihole.yourdomain.com
    http:
      paths:
      - backend:
          service:
            name: pihole
            port:
              number: 80
        path: /
        pathType: Prefix
```




