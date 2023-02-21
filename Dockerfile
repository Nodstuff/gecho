FROM golang:alpine AS builder

ARG TARGETOS
ARG TARGETARCH

WORKDIR $GOPATH/src/gecho/
COPY . .
RUN go get -d -v
RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -ldflags="-w -s" -o /go/bin/gecho
FROM scratch
COPY --from=builder /go/bin/gecho /go/bin/gecho
ENTRYPOINT ["/go/bin/gecho"]