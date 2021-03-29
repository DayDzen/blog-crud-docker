FROM golang:latest

ENV GO111MODULE on

WORKDIR /go/src/github.com/DayDzen/blog-crud-docker
COPY . .

RUN go get ./...
RUN go install ./...

EXPOSE 8000

CMD ["go", "run", "main.go"]