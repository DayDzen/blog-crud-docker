FROM golang:latest

ENV GO111MODULE on

WORKDIR /app
COPY . .

RUN go get github.com/cespare/reflex
RUN go get ./...
RUN go install ./...

ENTRYPOINT ["reflex", "-c", "reflex.conf"]