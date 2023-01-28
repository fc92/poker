#!/bin/sh

# start server in the background
/poker server -websocket localhost:8080 & 

# use tty2web to allow client to connect from a browser on port specified in $1
# --title-format will set a title for the browser window, with the possibility to specify a Room name in $2
# --permit-arguments allows to specify a player name in the URL
# -w allows to send mouse and keyboard instruction from the browser to the terminal on server side
#
# each client gets a dedicated terminal by calling /bin/sh
/tty2web --title-format "Poker $2" --permit-arguments -w -p $1 /bin/sh poker.sh