FROM golang:1.15-alpine
WORKDIR /app/
COPY ./bin/api-pipeline /app/
RUN ln -s /app/api-pipeline* /sbin/

