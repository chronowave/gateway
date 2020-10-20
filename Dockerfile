FROM golang:1.15.2-alpine3.12

COPY . /src
WORKDIR /src

RUN go build

ENTRYPOINT ["./gateway"]