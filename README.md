# PubSub Emulator UI

# Usage
```
docker run -it -p 8080:80 \
  -e PUBSUB_EMULATOR_HOST={HOST:PORT} \
  -e PUBSUB_PROJECT_ID=my-project-id \
  -e PUBSUB_TOPIC_ID=my-topic \
  -e PUBSUB_SUBSCRIPTION_ID=my-topic-subscription \
  yamato787/pubsub-emulator-ui
```

※ `{HOST}`はPubSub Emulatorの起動しているホストを指定してください。  
ホストマシン上で起動している場合は、`172.17.0.1:8085(環境によっては異なります)`や`host.docker.internal:8085(Mac/Windows)` を指定。  
同じDocker Network内でEmulatorを起動している場合はそのサービス名(`pubsub:8085`など)を指定してください。  

このリポジトリに同梱されている`docker-compose.yml`では、同一DockerネットワークでEmulatorを起動しているため、`pubsub:8085`を指定しています。
