FROM ubuntu:20.04

ARG DEBIAN_FRONTEND=noninteractive

RUN apt-get -y update &&\
    apt-get -y upgrade &&\
    apt -y install build-essential texinfo &&\
    apt-get -y install nano &&\
    apt-get -y install curl &&\
    apt-get -y install wget &&\
    cd usr/local &&\
    wget "https://go.dev/dl/go1.18.linux-amd64.tar.gz" &&\
    tar -xvf go1.18.linux-amd64.tar.gz &&\
    rm -rf /usr/local/go && tar -C /usr/local -xzf go1.18.linux-amd64.tar.gz

ENV PATH /usr/local/go/bin:$PATH

RUN mkdir go-gql-blogs && cd go-gql-blogs &&\
    go mod init github.com/nadirbasalamah/go-gql-blogs &&\
    go get github.com/99designs/gqlgen@v0.17.10

COPY . .

RUN mv server.go /go-gql-blogs &&\
    mv gqlgen.yml /go-gql-blogs &&\
    mv graph /go-gql-blogs &&\
    cd go-gql-blogs &&\ 
    go build

RUN curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b $(go env GOPATH)/bin 
