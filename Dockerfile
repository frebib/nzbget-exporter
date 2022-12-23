ARG EXPORTER_VER=0.2.2

FROM golang:alpine3.15

WORKDIR /tmp/enumerx
RUN apk add git && \
    git clone https://github.com/frebib/enumerx.git -b go1.18 . && \
    go install

WORKDIR /build
ADD go.mod go.sum ./
RUN go mod download

ARG EXPORTER_VER
ADD . ./
RUN go generate && \
    go build \
        -v \
        -ldflags="-w -s -X 'main.Version=$EXPORTER_VER'" \
        -o /nzbget_exporter

# ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

FROM spritsail/alpine:3.14

ARG EXPORTER_VER

LABEL org.opencontainers.image.authors="frebib <nzbget-exporter@frebib.net>" \
      org.opencontainers.image.title="nzbget-exporter" \
      org.opencontainers.image.url="https://github.com/frebib/nzbget-exporter" \
      org.opencontainers.image.description="NZBGet Prometheus metrics exporter" \
      org.opencontainers.image.version=${EXPORTER_VER}

COPY --from=0 /nzbget_exporter /usr/bin
CMD ["/usr/bin/nzbget_exporter"]
