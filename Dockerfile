FROM golang:latest as build

ADD . /build
WORKDIR /build
RUN CGO_ENABLED=0 go build -o /build/app.exe /build/main

FROM scratch

COPY --from=build /build/app.exe .

EXPOSE 8080

ENTRYPOINT ["/app.exe"]
