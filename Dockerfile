FROM golang

WORKDIR $GOPATH/src/godocker/main//
ADD . $GOPATH/src/godocker
RUN go build main.go
EXPOSE 8080

ENTRYPOINT ["./main"]