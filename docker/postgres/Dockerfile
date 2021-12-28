# syntax=docker/dockerfile:1
FROM mysql:8.0.27
MAINTAINER Arngrimur Bjarnason

ENV MYSQL_DATABASE csn_db
EXPOSE 3306
COPY data/create.sql /docker-entrypoint-initdb.d