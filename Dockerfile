FROM wardknight/wapti-test:latest

RUN apt-get update \ 
    && apt-get -y install nmap

COPY server /usr/bin/server

EXPOSE 18080

CMD [ "server" ]