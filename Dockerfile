FROM golang:alpine

WORKDIR /slackhell
COPY . /slackhell
RUN go build -o /usr/bin/slackhell cmd/main.go

CMD ["/usr/bin/slackhell","run"]