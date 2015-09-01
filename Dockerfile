FROM golang:1.5

WORKDIR /go/src/bitbucket.org/emicklei/v8dispatcher
ADD . /go/src/bitbucket.org/emicklei/v8dispatcher

CMD make dockerbuild