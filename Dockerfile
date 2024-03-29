FROM golang:1.20 AS builder

WORKDIR /go/src/app
COPY . .

RUN CGO_ENABLED=0 go build -o ./ ./cmd/...

FROM alpine:3

COPY --from=builder /go/src/app/kube-plex /go/src/app/transcode-launcher /