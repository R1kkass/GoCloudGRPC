FROM golang:1.21 as base

FROM base as dev

RUN curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b $(go env GOPATH)/bin

RUN go install github.com/go-delve/delve/cmd/dlv@latest
RUN export GOPATH=/go

WORKDIR /var/www/html

CMD air -c .air.toml -- -h
#  ; dlv exec --continue --accept-multiclient --listen=:2345 --headless=true --api-version=2 --log --log-output=rpc,dap,debugger ./tmp/main
