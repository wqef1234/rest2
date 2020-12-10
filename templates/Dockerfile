
FROM golang:1.12.0-alpine3.9

RUN apk add --no-cache git

RUN mkdir /rest2

ADD . /rest2

WORKDIR /rest2

RUN go mod download

RUN go build -o main .

CMD ["/resttest/main"]