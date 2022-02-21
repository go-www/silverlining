FROM golang:alpine

ADD . /build
WORKDIR /build
RUN CGO_ENABLED=0 go build -o /build/app.exe /build/main

ENTRYPOINT ["/build/app.exe"]
