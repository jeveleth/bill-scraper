# build the binary
FROM golang:alpine AS builder

RUN apk --no-cache add ca-certificates

WORKDIR /go/src/app
COPY . .
RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -v -o app

# build a small image that runs the binary
FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

COPY --from=builder /go/src/app .
CMD ["./app"]