# Creates a test DB for local development.
FROM mysql:latest AS db_builder
ARG PASS
ENV MYSQL_ROOT_PASSWORD=$PASS
ENV MYSQL_DATABASE=pod_tracker
RUN mkdir conf.d/
ADD custom.cnf /conf.d
EXPOSE 3306
