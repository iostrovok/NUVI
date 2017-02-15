# vim: set softtabstop=2 shiftwidth=2:
# redis-server /usr/local/etc/redis.conf

SHELL = bash


CURDIR = ${PWD}
export GOPATH=${CURDIR}

all: clean install start

clean:
	@echo "Clean server" 
	$(shell rm ./bin/index)


install:
	@echo "Build server"
	go get github.com/garyburd/redigo/redis
	go build -o ./bin/index ./src/index.go


start:
	@echo "Start server" 
	./bin/index -dir ${CURDIR} -port 19720

test: test-manager test-extract-urls test-loader test-redis

test-manager:
	@echo "test manager"
	go get gopkg.in/check.v1 
	go test ./src/manager/

test-loader:
	@echo "test loader"
	go get gopkg.in/check.v1 
	go test ./src/loader/


test-extract-urls:
	@echo "test extract_urls"
	go get gopkg.in/check.v1 
	go test ./src/extract_urls/

test-redis:
	@echo "test redis"
	go get gopkg.in/check.v1 
	go get github.com/garyburd/redigo/redis
	go test ./src/redis/







