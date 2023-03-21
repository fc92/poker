#!/bin/sh

DEBUG_PORT=$1
SERVER=localhost:8080
# start server in the background
/poker server -websocket $SERVER & 

/dlv debug /go/src/cmd -r stdin:/tmp/stdin.txt --headless --listen=:$DEBUG_PORT --api-version=2 -- client -websocket $SERVER

# use kubect attach -it to access the UI