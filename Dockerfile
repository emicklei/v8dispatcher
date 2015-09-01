FROM golang:1.5

# to upgrade git
#RUN apt-get update
#RUN apt-get upgrade
#RUN apt-get install software-properties-common python-software-properties -y
#RUN add-apt-repository ppa:git-core/ppa

RUN apt-get update
RUN apt-get upgrade
RUN apt-get -y install git make g++ && apt-get clean
RUN git --version
RUN make --version
RUN g++ --version

# depot tools
RUN git clone https://chromium.googlesource.com/chromium/tools/depot_tools.git /usr/local/depot_tools
ENV PATH $PATH:/usr/local/depot_tools

# v8worker
RUN git clone https://github.com/ry/v8worker.git /go/src/github.com/ry/v8worker
WORKDIR /go/src/github.com/ry/v8worker
RUN make && make install

WORKDIR /go/src/bitbucket.org/emicklei/v8dispatcher
ADD . /go/src/bitbucket.org/emicklei/v8dispatcher

CMD make dockerbuild