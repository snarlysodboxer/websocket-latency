#!/bin/bash
set -e

docker pull snarlysodboxer/websocket-latency:builder
docker pull snarlysodboxer/websocket-latency:latest

docker build --target builder \
       --cache-from=snarlysodboxer/websocket-latency:builder \
       --tag snarlysodboxer/websocket-latency:builder .

docker build --target runtime \
       --cache-from=snarlysodboxer/websocket-latency:builder \
       --cache-from=snarlysodboxer/websocket-latency:latest \
       --tag snarlysodboxer/websocket-latency:latest .

docker tag snarlysodboxer/websocket-latency:latest snarlysodboxer/websocket-latency:0.0.0

docker push snarlysodboxer/websocket-latency:builder
docker push snarlysodboxer/websocket-latency:latest
docker push snarlysodboxer/websocket-latency:0.0.0
