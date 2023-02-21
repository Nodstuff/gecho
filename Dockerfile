FROM golang:alpine AS builder
WORKDIR $GOPATH/src/gecho/
COPY . .
RUN go get -d -v
RUN GOOS=linux GOARCH=arm64 go build -ldflags="-w -s" -o /go/bin/gecho
FROM scratch
COPY --from=builder /go/bin/gecho /go/bin/gecho
ENTRYPOINT ["/go/bin/gecho"]