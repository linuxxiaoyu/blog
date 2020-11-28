FROM golang:latest

ENV GOPROXY https://goproxy.cn,direct
WORKDIR /go/src/app
COPY . .

RUN go build .
EXPOSE 8080
ENTRYPOINT [ "./blog" ]