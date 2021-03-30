FROM golang:1.14.4 AS builder
WORKDIR /go/src/github.com/snarlysodboxer/websocket-latency
ADD main.go go.mod go.sum /go/src/github.com/snarlysodboxer/websocket-latency/
RUN go build main.go

FROM scratch AS runtime
COPY --from=builder /go/src/github.com/snarlysodboxer/websocket-latency /websocket-latency
