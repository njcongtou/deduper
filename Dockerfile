FROM golang

WORKDIR $GOPATH/src/godocker/main//
ADD . $GOPATH/src/godocker
RUN go build main.go
EXPOSE 8001 8002 8003 9999

ENTRYPOINT ["./main"]