FROM golang:1.12rc1-alpine3.8

RUN go build main.go

EXPOSE 81

CMD ["main.exe"]