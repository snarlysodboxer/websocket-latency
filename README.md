# websocket-latency

## Simulate latency for websockets

Listen on an address, forward each websocket message to a destination address after a delay.
Messages are delayed in both directions.

```shell
docker run -it --rm snarlysodboxer:0.0.0 -listenAddress 0.0.0.0:9090 -forwardAddress localhost:8080 -delay 300ms
```
