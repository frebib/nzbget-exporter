ARG EXPORTER_VER=0.2.1

FROM golang:alpine3.14

RUN apk add git && \
    go install github.com/niktri/enumerx@latest

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

LABEL maintainer="frebib <nzbget-exporter@frebib.net>" \
      org.label-schema.vendor="frebib" \
      org.label-schema.name="nzbget-exporter" \
      org.label-schema.url="https://github.com/frebib/nzbget-exporter" \
      org.label-schema.description="NZBGet Prometheus metrics exporter" \
      org.label-schema.version=${EXPORTER_VER}

COPY --from=0 /nzbget_exporter /usr/bin
CMD ["/usr/bin/nzbget_exporter"]
