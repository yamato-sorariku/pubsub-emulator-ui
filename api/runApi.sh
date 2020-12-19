#!/bin/sh

export PUBSUB_EMULATOR_HOST=localhost:8085
export PUBSUB_PROJECT_ID=my-project-id
export PUBSUB_TOPIC_ID=my-topic
export PUBSUB_SUBSCRIPTION_ID=my-topic-subscription

go run main.go
