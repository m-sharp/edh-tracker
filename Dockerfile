FROM golang:1.20.0-alpine3.17 as TrackerAPI

ENV DBHOST ""
ENV DBUSER ""
ENV DBPASSWORD ""
ENV DBPORT ""
ENV CSRFSEC ""

WORKDIR /edh-tracker
RUN mkdir build/

COPY vendor/ vendor/
COPY go.mod .
COPY go.sum .

COPY api.go .
COPY main.go .
COPY lib/ lib/

RUN go build -o build/ .
CMD ["build/edh-tracker"]
