FROM ubuntu

COPY bin/appdaemon.linux /appdaemon

EXPOSE 6380

CMD ["/appdaemon", "-c", "/config.yaml", "-disable-daemon"]
