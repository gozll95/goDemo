FROM golang

RUN mkdir -p /go/src/go-kit-demo1

ADD . /go/src/go-kit-demo1

RUN go get  -t -v ./...
RUN go get  github.com/canthefason/go-watcher
RUN go install github.com/canthefason/go-watcher/cmd/watcher

ENTRYPOINT  watcher -run go-kit-demo1/notificator/cmd  -watch go-kit-demo1/notificator
