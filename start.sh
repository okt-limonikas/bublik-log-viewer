#!/usr/bin/env bash

start_server() {
    python3 -m http.server 4200 --bind 127.0.0.1 & # start in backgorund
    SERVER_PID=$!
}

stop_server() {
    if [[ -n "$SERVER_PID" ]]; then
        kill -SIGTERM "$SERVER_PID"
        wait "$SERVER_PID"
    fi
}

start_server

# Trap the SIGINT signal (Ctrl+C) to stop the server gracefully
trap stop_server SIGINT

sleep 2

echo "You can browse logs at http://localhost:4200"

wait "$SERVER_PID"
