#!/bin/sh

set -x
# This script is meant to be used with kubectl exec on a sleeping container
# (avoid kubectl attach to dlv process as stdin is used by dlv)
#
# debug server with : kubectl exec -it room-teamred-xxxxx -- /clients-debug.sh server 9229
# or
# debug client with : kubectl exec -it room-teamred-xxxxx -- /clients-debug.sh client 9229

TYPE=$1 # client or server depending on what needs to be debugged
DEBUG_PORT=$2
SERVER=localhost:8080
# start server in the background to debug client
if [[ "$TYPE" = "client" ]]
then
    echo Start server without debug
    /poker server -websocket $SERVER & 
fi

echo Start $TYPE with debug port $DEBUG_PORT
/dlv debug /go/src/cmd --headless --listen=:$DEBUG_PORT --api-version=2 -- $TYPE -websocket $SERVER
