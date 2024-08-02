FROM golang:1.21 as base

FROM base as dev

RUN curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b $(go env GOPATH)/bin

# RUN mkdir /var/www
# RUN mkdir /var/www/html

# COPY . /var/www/html

WORKDIR /var/www/html
CMD air --build.cmd="go build -buildvcs=false -o ./tmp/main ." --build.bin "./tmp/main"