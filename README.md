# edh-tracker
Magic the Gathering EDH Tracking for a small pod of players.

## Install

- Docker
- Go

## Dev Reference

- Pull down Go dependencies with `go mod vendor`
- Build Docker images:
  - `docker build -t registry.digitalocean.com/harp-do-registry/edh-tracker .`
  - `docker build --build-arg PASS=REDACTED -t edh-tracker-db ./mysql/`
- Run docker images:
  - ```
    docker run -p 8080:8081 -it \
      --env DBHOST=host.docker.internal \
      --env DBUSER=root \ 
      --env DBPASSWORD=REDACTED \
      --env DBPORT=3306 \
      --env DEV=1 \
      --env CSRFSEC=REDACTED registry.digitalocean.com/harp-do-registry/edh-tracker
    ```
  - `docker run --detach --name=edh-tracker-db --publish 3306:3306 edh-tracker-db`

## Required Environment Variables

- `DBHOST` - Hostname of database.
- `DBUSER` - Username to connect to the database with.
- `DBPASSWORD` - Password to connect to the database with.
- `DBPORT` - Port to connect to database on.
- `CSRFSEC` - Secret key for generating CSRF tokens. Should be a random, improbable to guess 32-byte long string.
