#!/bin/sh

case "$1" in
  "api_gateway")
    ./api_gateway.o
    ;;
  "subscriber")
    ./subscriber.o
    ;;
  *)
    echo "Invalid binary selection: $1"
    exit 1
    ;;
esac
