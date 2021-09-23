FROM golang:1.17 as build

WORKDIR /go/src/github.com/webdevops/deadmanssnitch-exporter

# Get deps (cached)
COPY ./go.mod /go/src/github.com/webdevops/deadmanssnitch-exporter
COPY ./go.sum /go/src/github.com/webdevops/deadmanssnitch-exporter
COPY ./Makefile /go/src/github.com/webdevops/deadmanssnitch-exporter
RUN make dependencies

# Compile
COPY ./ /go/src/github.com/webdevops/deadmanssnitch-exporter
RUN make test
RUN make lint
RUN make build
RUN ./deadmanssnitch-exporter --help

#############################################
# FINAL IMAGE
#############################################
FROM gcr.io/distroless/static
ENV LOG_JSON=1
COPY --from=build /go/src/github.com/webdevops/deadmanssnitch-exporter/deadmanssnitch-exporter /
USER 1000
ENTRYPOINT ["/deadmanssnitch-exporter"]
