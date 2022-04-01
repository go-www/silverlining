FROM golang:alpine

ADD . /build
WORKDIR /build
RUN CGO_ENABLED=0 go build -o /build/app.exe /build/main

RUN mv /build/app.exe /app.exe

# Cleanup
RUN rm -rf /build
RUN rm -rf /go
RUN rm -rf /root/.cache
RUN rm -rf $(go env GOPATH)
RUN rm -rf $(go env GOROOT)

ENTRYPOINT ["/app.exe"]
