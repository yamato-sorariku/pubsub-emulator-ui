FROM golang:1.15.5-buster as GoBuild
WORKDIR /app
COPY ./api /app
RUN go build

FROM node:14.15.0-buster as NodeBuild
WORKDIR /tmp
COPY ./front/package.json /tmp
COPY ./front/yarn.lock /tmp
RUN yarn install
WORKDIR /app
COPY ./front /app
RUN yarn install
RUN yarn generate

FROM nginx:1.18.0-alpine
COPY --from=NodeBuild /app/dist /usr/share/nginx/html
WORKDIR /app
COPY --from=GoBuild /app/pubsub-emulator-ui  /app
COPY docker/pubsub-emulator-ui/run.sh /run.sh
RUN chmod 700 /run.sh && mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
COPY docker/pubsub-emulator-ui/default.conf /etc/nginx/conf.d/default.conf

CMD /run.sh
