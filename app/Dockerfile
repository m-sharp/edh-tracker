# This Dockerfile should be run from the edh-tracker/ directory via the -f flag on docker build
FROM node:18-alpine AS reactBuild

# ToDo: Might be able to speed up builds by copying https://github.com/vercel/next.js/blob/canary/examples/with-docker/Dockerfile a bit more

RUN apk add --no-cache libc6-compat
WORKDIR /app

COPY app/package.json .
COPY app/package-lock.json .
COPY app/next.config.js .
COPY app/public/ public/
COPY app/src/ src/

RUN npm ci
RUN npm run build

FROM golang:1.20.0-alpine3.17 as trackerAPP

WORKDIR /edh-tracker
RUN mkdir build/

COPY vendor/ vendor/
COPY go.mod .
COPY go.sum .

COPY app/main.go .
COPY lib/ lib/
COPY --from=reactBuild /app/out/ app/

RUN go build -o build/app .
CMD ["build/app"]
