
ARG GO_VERSION=1.16
FROM golang:${GO_VERSION}-alpine AS builder

ARG VERSION=""
ARG BRANCH=""
ARG COMMIT=""

RUN mkdir /user && \
    echo 'nobody:x:65534:65534:nobody:/:' > /user/passwd && \
    echo 'nobody:x:65534:' > /user/group

WORKDIR /src
COPY ./go.mod ./go.sum ./
RUN go mod download

COPY ./ ./ 

# Build the executable to `/app`. Mark the build as statically linked.
RUN CGO_ENABLED=0 go build \
    -installsuffix 'static' \
    -ldflags "-X main.Revision=${COMMIT} \
        -X main.Branch=${BRANCH} \
        -X main.Version=${VERSION} \
        -X main.BuildDate=$(date -u ""+%Y%m%d-%H:%M:%S"")" \
    -o /app .


FROM scratch AS final
LABEL maintainer="Kris@budd.ee"
COPY --from=builder /user/group /user/passwd /etc/
COPY --from=builder /app /app
EXPOSE 80
USER nobody:nobody
ENTRYPOINT ["/app"]