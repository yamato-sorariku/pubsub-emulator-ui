#!/bin/sh

# Run API
export GIN_MODE=release
/app/pubsub-emulator-ui &

# Nginx
nginx -g "daemon off;"
