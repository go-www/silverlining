FROM golang:alpine as builder

ADD . /build
WORKDIR /build
RUN CGO_ENABLED=0 go build -o /build/app.exe /build/main

FROM scratch

COPY --from=builder /build/app.exe /app.exe

ENTRYPOINT ["/app.exe"]
