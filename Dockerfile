FROM debian:stable-slim

WORKDIR /app
COPY go-emqx-to-tdengine config.toml create.sql insert.sql /app/

CMD ["go-emqx-to-tdengine"]