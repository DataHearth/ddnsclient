FROM golang:1.15 as base

WORKDIR /go/src/ddnsclient
ADD . /go/src/ddnsclient

RUN go get -d -v ./...

RUN go build -o /go/bin/ddnsclient cmd/main.go

FROM gcr.io/distroless/base

COPY --from=base /go/bin/ddnsclient /

CMD ["/ddnsclient"]