FROM golang:1.12rc1-alpine3.8

COPY main.go .

RUN go build main.go

EXPOSE 81

CMD ["main"]