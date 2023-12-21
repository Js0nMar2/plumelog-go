FROM debian:10

COPY config.toml config.toml
ADD app app
ADD dist /dist
RUN chmod -R 755 app
ENTRYPOINT ["./app"]
