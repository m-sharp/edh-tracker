# edh-tracker
Magic the Gathering EDH Tracking for a pod of players using Go, MySQL, & React.

## Install

- Docker
- Go
- Node.js

## Dev Reference

- Pull down Go dependencies with `go mod vendor`
- Pull down npm dependencies with `cd app && npm install`
- Build Docker images:
  - `docker build --build-arg PASS=REDACTED -t edh-tracker-db ./mysql/`
  - `docker build -t registry.digitalocean.com/harp-do-registry/edh-tracker .`
  - `docker build -f app/Dockerfile -t registry.digitalocean.com/harp-do-registry/edh-tracker-app .`
- Run docker images:
  - Run DB: `docker run --detach --name=edh-tracker-db --publish 3306:3306 edh-tracker-db`
  - Run API server:
    ```shell
    docker run -p 8080:8081 -it \
      --env DBHOST=host.docker.internal \
      --env DBUSER=root \ 
      --env DBPASSWORD=REDACTED \
      --env DBPORT=3306 \
      --env DEV=1 registry.digitalocean.com/harp-do-registry/edh-tracker
    ```
  - Run React web app: `docker run -p 8081:8081 -it registry.digitalocean.com/harp-do-registry/edh-tracker-app:latest`

## Required Environment Variables

- API Server:
  - `DBHOST` - Hostname of database.
  - `DBUSER` - Username to connect to the database with.
  - `DBPASSWORD` - Password to connect to the database with.
  - `DBPORT` - Port to connect to database on.
