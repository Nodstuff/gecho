FROM golang:alpine AS builder

ARG TARGETOS
ARG TARGETARCH

WORKDIR $GOPATH/src/gecho/
COPY . .

RUN go get -d -v
RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -ldflags="-w -s" -o /go/bin/gecho
RUN mkdir -p /ssl/certs

FROM scratch
COPY --from=builder /go/bin/gecho /go/bin/gecho
COPY --from=builder /ssl/certs /ssl/certs

ENTRYPOINT ["/go/bin/gecho"]