FROM golang:1.13

WORKDIR /go/src/crudapp
COPY . /go/src/crudapp
RUN go mod download
RUN make build

EXPOSE 8080/tcp

CMD [ "/go/src/crudapp/bin/crudapp" ]