FROM golang:1.26.0-alpine3.23 AS tracker_api

ENV DBHOST=""
ENV DBUSER=""
ENV DBPASSWORD=""
ENV DBPORT=""
ENV CSRFSEC=""

WORKDIR /edh-tracker
RUN mkdir build/

COPY vendor/ vendor/
COPY go.mod .
COPY go.sum .

COPY api.go .
COPY main.go .
COPY lib/ lib/
COPY data/ data/

RUN go build -o build/ .
CMD ["build/edh-tracker"]
