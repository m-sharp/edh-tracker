# This Dockerfile should be run from the edh-tracker/ directory via the -f flag on docker build
# ToDo: Might be able to speed up builds by copying https://github.com/vercel/next.js/blob/canary/examples/with-docker/Dockerfile a bit more
FROM node:18-alpine AS reactBuild

RUN apk add --no-cache libc6-compat
WORKDIR /app

# Copy JS dependencies
COPY app/package.json .
COPY app/package-lock.json .
COPY app/tsconfig.json .

# Build JS env
RUN npm ci

# Copy app source files
COPY app/public/ public/
COPY app/src/ src/

# Build app
RUN npm run build

FROM golang:1.20.0-alpine3.17 as goBuild

WORKDIR /go-build
RUN mkdir build/

# Copy Go dependencies
COPY vendor/ vendor/
COPY go.mod .
COPY go.sum .

# Copy Go Source
COPY app/main.go .
COPY lib/ lib/

# Build go webserver
RUN go build -o build/web-server .

FROM golang:1.20.0-alpine3.17 as app

WORKDIR /edh-tracker

COPY --from=goBuild /go-build/build/web-server .
COPY --from=reactBuild /app/build/ app/

CMD ["./web-server"]
