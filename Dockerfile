FROM sammascanner/wapiti:latest

RUN apt-get update \ 
    && apt-get -y install nmap

COPY server /usr/bin/server

EXPOSE 8080

CMD [ "server" ]