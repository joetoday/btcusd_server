FROM golang:latest

LABEL maintainer="Tunde Ogundele <ogundele.tj@gmail.com>"

RUN mkdir /btcusd_server

WORKDIR /btcusd_server
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

RUN go build -o main .

EXPOSE 80

CMD ["/btcusd_server/main"]

