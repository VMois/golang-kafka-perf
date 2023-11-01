FROM golang:1.21.3

RUN mkdir /opt/app
WORKDIR /opt/app

ADD ./scripts/entrypoint.sh entrypoint.sh
RUN chmod +x entrypoint.sh
COPY ./api_gateway api_gateway
COPY ./subscriber subscriber

RUN cd /opt/app/subscriber && go mod tidy && go build -o ../subscriber.o && chmod +x ../subscriber.o
RUN cd /opt/app/api_gateway && go mod tidy && go build -o ../api_gateway.o && chmod +x ../api_gateway.o

ENTRYPOINT ["/opt/app/entrypoint.sh"]
