FROM node:18-alpine AS reactBuild

RUN apk add --no-cache libc6-compat
WORKDIR /app

COPY app/package.json app/package-lock.json* ./
COPY app/next.config.js ./
COPY app/public/ public/
COPY app/src/ src/
RUN npm ci
RUN npm run build

FROM golang:1.20.0-alpine3.17 as tracker

ENV DBHOST ""
ENV DBUSER ""
ENV DBPASSWORD ""
ENV DBPORT ""
ENV CSRFSEC ""

RUN mkdir /edh-tracker
WORKDIR /edh-tracker
RUN mkdir build/

COPY vendor/ vendor/
COPY go.mod .
COPY go.sum .

COPY main.go .
COPY lib/ lib/
COPY web/ web/
COPY --from=reactBuild /app/out/ app/

RUN go build -o build/ ./...
CMD ["build/edh-tracker"]
