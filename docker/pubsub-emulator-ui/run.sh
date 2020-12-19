#!/bin/sh

# Run API
/app/pubsub-emulator-ui &

# Nginx
nginx -g "daemon off;"
