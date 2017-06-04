FROM alpine:latest

RUN adduser -D -g '' -u 1000 app-user

RUN install -o app-user -g app-user -m 0755 -d /app

RUN apk --update upgrade && \
  apk add --no-cache ca-certificates && \
  update-ca-certificates && \
  rm -rf /var/cache/apk/*

USER app-user

ENV PATH $PATH:/app

WORKDIR /app

ADD ./tor-status-proxy /app/

CMD ["tor-status-proxy"]

