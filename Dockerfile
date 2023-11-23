FROM golang:1.20.0-alpine3.17 as tracker

ENV DBHOST ""
ENV DBUSER ""
ENV DBPASSWORD ""
ENV DBPORT ""
ENV CSRFSEC ""

RUN mkdir /edh-tracker
WORKDIR /edh-tracker
RUN mkdir app/

COPY web/ web/

COPY go.mod .
COPY go.sum .
COPY main.go .
COPY lib/ lib/
COPY vendor/ vendor/

RUN go build -o app/ ./...
CMD ["app/edh-tracker"]
